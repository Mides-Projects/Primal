package service

import (
    "context"
    "errors"
    quark "github.com/Mides-Projects/Quark"
    "github.com/holypvp/primal/common"
    "github.com/holypvp/primal/model"
    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "strings"
    "sync"
    "time"
)

type AccountService struct {
    col *mongo.Collection

    accountsMu sync.RWMutex
    accounts   map[string]*model.Account

    ttlCacheMu sync.RWMutex
    ttlCache   *quark.Quark[*model.Account]

    accountsIdMu sync.RWMutex
    accountsId   map[string]string
}

// LookupById retrieves an account by its ID. It's safe to use this method because
// it's only reading the map.
func (s *AccountService) LookupById(id string) *model.Account {
    s.accountsMu.RLock()
    defer s.accountsMu.RUnlock()

    if acc, ok := s.accounts[id]; ok {
        return acc
    }

    s.ttlCacheMu.RLock()
    defer s.ttlCacheMu.RUnlock()

    if acc, ok := s.ttlCache.Get(id); ok {
        return acc
    }

    return nil
}

// LookupByName retrieves an account by its name. It's safe to use this method because
// it's only reading the map.
func (s *AccountService) LookupByName(name string) *model.Account {
    s.accountsIdMu.RLock()
    defer s.accountsIdMu.RUnlock()

    if id, ok := s.accountsId[strings.ToLower(name)]; ok {
        return s.LookupById(id)
    }

    return nil
}

// UnsafeLookupById retrieves an account by its ID. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *AccountService) UnsafeLookupById(id string, keep bool) (*model.Account, error) {
    if acc := s.LookupById(id); acc != nil {
        return acc, nil
    }

    // These values are only cached for 72 hours, after that they are removed from redis
    // but still available in our mongo database.
    // If the account was fetch from database, it will be cached into redis to prevent further database calls in the next 72 hours.
    var (
        acc *model.Account
        err error
    )
    if acc, err = s.lookupAtRedis("ids:", id); err != nil {
        return nil, err
    } else if acc == nil {
        acc, err = s.lookupAtMongo("_id", id)
    }

    if err != nil {
        return nil, err
    } else if acc == nil {
        return nil, nil
    }

    s.Cache(acc, keep)

    return acc, nil
}

// UnsafeLookupByName retrieves an account by its name. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *AccountService) UnsafeLookupByName(name string, keep bool) (*model.Account, error) {
    if acc := s.LookupByName(name); acc != nil {
        return acc, nil
    }

    var (
        acc *model.Account
        err error
    )
    if acc, err = s.lookupAtRedis("names:", name); err != nil {
        return nil, err
    } else if acc == nil {
        acc, err = s.lookupAtMongo("name", name)
    }

    if err != nil {
        return nil, err
    } else if acc == nil {
        return nil, nil
    }

    s.Cache(acc, keep)

    return acc, nil
}

// UpdateName updates the name of an account.
func (s *AccountService) UpdateName(oldName, newName, id string) {
    s.accountsIdMu.Lock()

    delete(s.accountsId, strings.ToLower(oldName))
    s.accountsId[strings.ToLower(newName)] = id

    s.accountsIdMu.Unlock()
}

// Cache caches an account.
func (s *AccountService) Cache(a *model.Account, keep bool) {
    if s.ttlCache == nil {
        panic("service not hooked to the ttlCache")
    }

    if keep {
        s.accountsMu.Lock()
        s.accounts[a.Id()] = a
        s.accountsMu.Unlock()
    } else {
        s.ttlCacheMu.Lock()
        s.ttlCache.Set(a.Id(), a)
        s.ttlCacheMu.Unlock()
    }

    s.accountsIdMu.Lock()
    s.accountsId[strings.ToLower(a.Name())] = a.Id()
    s.accountsIdMu.Unlock()
}

// Invalidate invalidates an account.
func (s *AccountService) Invalidate(acc *model.Account) {
    s.ttlCacheMu.Lock()
    s.ttlCache.Invalidate(acc.Id())
    s.ttlCacheMu.Unlock()

    s.accountsMu.Lock()
    delete(s.accounts, acc.Id())
    s.accountsMu.Unlock()

    s.accountsIdMu.Lock()
    delete(s.accountsId, strings.ToLower(acc.Name()))
    s.accountsIdMu.Unlock()
}

// RedisCache caches an account into the Redis database.
func (s *AccountService) RedisCache(acc *model.Account) error {
    if common.RedisClient == nil {
        return errors.New("redis client not found")
    }

    pip := common.RedisClient.Pipeline()
    if pip == nil {
        return errors.New("redis pipeline not found")
    }

    pip.Set(context.Background(), "primal%ids:"+acc.Id(), acc.MarshalString(), 72*time.Hour)
    pip.Set(context.Background(), "primal%names:"+strings.ToLower(acc.Name()), acc.MarshalString(), 72*time.Hour)

    if _, err := pip.Exec(context.Background()); err != nil {
        return err
    }

    return nil
}

// Update updates an account.
func (s *AccountService) Update(acc *model.Account) error {
    if s.col == nil {
        return errors.New("service not hooked to the database")
    }

    res, err := s.col.UpdateOne(
        context.TODO(),
        bson.D{{"_id", acc.Id()}},
        bson.D{{"$set", acc.Marshal()}},
        options.Update().SetUpsert(true),
    )
    if err != nil {
        return err
    }

    if res.UpsertedCount > 0 {
        common.Log.Printf("Account %s was inserted", acc.Id())
    } else {
        common.Log.Printf("Account %s was updated", acc.Id())
    }

    if err = s.RedisCache(acc); err != nil {
        return err
    }

    return nil
}

func (s *AccountService) lookupAtRedis(k, v string) (*model.Account, error) {
    val, err := common.RedisClient.Get(context.Background(), "primal%"+k+v).Result()
    if errors.Is(err, redis.Nil) {
        return nil, nil
    } else if err != nil {
        return nil, err
    } else if val == "" {
        return nil, errors.New("empty value")
    } else {
        acc := &model.Account{}
        if acc.UnmarshalString(val) != nil {
            return nil, err
        }

        return acc, nil
    }
}

func (s *AccountService) lookupAtMongo(k, v string) (*model.Account, error) {
    var acc model.Account

    if err := s.col.FindOne(context.TODO(), bson.D{{k, v}}).Decode(&acc); err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }

        return nil, err
    }

    return &acc, nil
}

// Hook hooks the account service to the database.
func (s *AccountService) Hook(db *mongo.Database) error {
    if s.col != nil {
        return errors.New("service already hooked to the database")
    }

    s.col = db.Collection("trackers")

    s.ttlCache = quark.New[*model.Account](2*time.Hour, 2*time.Hour)
    s.ttlCache.SetListener(func(key string, value *model.Account, reason quark.Reason) {
        if reason == quark.ManualReason {
            return
        }

        s.accountsIdMu.Lock()
        delete(s.accountsId, strings.ToLower(value.Name()))
        s.accountsIdMu.Unlock()
    })

    return nil
}

// Account returns the account service.
func Account() *AccountService {
    return accountService
}

var accountService = &AccountService{
    accounts:   make(map[string]*model.Account),
    accountsId: make(map[string]string),
}
