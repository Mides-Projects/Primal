package redis

import (
    "context"
    "errors"
    "github.com/bytedance/sonic"
    "github.com/holypvp/primal/common"
    "github.com/holypvp/primal/protocol"
    "github.com/redis/go-redis/v9"
    "time"
)

type Service struct {
    ctx context.Context

    key string
}

func NewService(key string) Service {
    return Service{
        ctx: context.Background(),
        key: key,
    }
}

// LookupJSON retrieves a value from the Redis storage.
func (s Service) LookupJSON(key string) (map[string]interface{}, error) {
    client := common.RedisClient
    if client == nil {
        return nil, errors.New("redis client is not initialized")
    }

    encoded, err := client.Get(s.ctx, s.key+key).Result()
    if err != nil {
        if errors.Is(err, redis.Nil) {
            return nil, nil
        }

        return nil, err
    }

    var decoded map[string]interface{}
    if err = sonic.Unmarshal([]byte(encoded), &decoded); err != nil {
        return nil, err
    }

    return decoded, nil
}

// LookupString retrieves a value from the Redis storage.
func (s Service) LookupString(key string) (string, error) {
    client := common.RedisClient
    if client == nil {
        return "", errors.New("redis client is not initialized")
    }

    res, err := client.Get(s.ctx, s.key+key).Result()
    if err != nil {
        if errors.Is(err, redis.Nil) {
            return "", nil
        }

        return "", err
    }

    return res, nil
}

// StoreJSON stores a value into the Redis storage.
func (s Service) StoreJSON(key string, value map[string]interface{}, ttl time.Duration) error {
    client := common.RedisClient
    if client == nil {
        return errors.New("redis client is not initialized")
    }

    encoded, err := sonic.Marshal(value)
    if err != nil {
        return err
    }

    return client.Set(s.ctx, s.key+key, encoded, ttl).Err()
}

// StoreString stores a value into the Redis storage.
func (s Service) StoreString(key, value string, ttl time.Duration) error {
    client := common.RedisClient
    if client == nil {
        return errors.New("redis client is not initialized")
    }

    return client.Set(s.ctx, s.key+key, value, ttl).Err()
}

func Publish(packet protocol.Packet) {
    client := common.RedisClient
    if client == nil {
        return
    }

    buf := protocol.NewWriter()

    pid := packet.ShieldId()
    buf.Varint32(&pid)
    packet.Marshal(buf)

    if err := client.Publish(context.Background(), "primal", buf).Err(); err != nil {
        common.Log.Printf("failed to publish packet: %v", err)
    }
}
