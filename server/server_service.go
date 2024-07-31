package server

import (
    "context"
    "github.com/holypvp/primal/server/model"
    "github.com/holypvp/primal/server/object"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

var (
    serversMap = map[string]*object.ServerInfo{}

    groupsMap = map[string]*object.ServerGroup{}

    collectionServers *mongo.Collection
    collectionGroups  *mongo.Collection
)

type ServerService struct{}

func SaveModel(infoModel model.ServerInfoModel) {
    if collectionServers == nil {
        log.Fatal("Servers collection is nil")
    }

    result, err := collectionServers.UpdateOne(
        context.TODO(),
        bson.D{{"_id", infoModel.Id}},
        bson.D{{"$set", infoModel}},
        options.Update().SetUpsert(true),
    )
    // _, err := collectionServers.InsertOne(context.TODO(), infoModel)
    if err != nil {
        log.Fatal(err)

        return
    }

    if result.UpsertedCount > 0 {
        log.Printf("Server %s was inserted", infoModel.Id)
    } else {
        log.Printf("Server %s was updated", infoModel.Id)
    }
}

func (service *ServerService) LoadGroups(database *mongo.Database) {
    if collectionGroups != nil {
        log.Fatal("Groups collection is already set")
    }

    collectionGroups = database.Collection("serverGroups")
    cursor, err := collectionGroups.Find(context.TODO(), bson.D{{}})
    if err != nil {
        panic(err)
    }

    for cursor.Next(context.Background()) {
        var result = &model.ServerGroupModel{}
        if err := cursor.Decode(result); err != nil {
            panic(err)
        }

        // TODO: Implement this
    }

    if err := cursor.Err(); err != nil {
        panic(err)

        return
    }

    if err = cursor.Close(context.TODO()); err != nil {
        panic(err)
    }

    log.Printf("Successfully loaded %d server groups", len(groupsMap))
}

func (service *ServerService) LoadServers(database *mongo.Database) {
    collectionServers = database.Collection("servers")
    if collectionServers == nil {
        log.Fatal("Servers collection is nil")
    }

    cursor, err := collectionServers.Find(context.TODO(), bson.D{{}})
    if err != nil {
        panic(err)
    }

    for cursor.Next(context.Background()) {
        var result = &model.ServerInfoModel{}

        if err = cursor.Decode(result); err != nil {
            log.Println("Failed to decode server info model: ", err)

            continue
        }

        service.StoreServer(&object.ServerInfo{
            id:     result.Id,
            port:   result.Port,
            groups: result.Groups,

            maxSlots:    result.MaxSlots,
            heartbeat:   result.Heartbeat,
            bungeeCord:  result.BungeeCord,
            onlineMode:  result.OnlineMode,
            initialTime: result.InitialTime,
        })
    }

    if err := cursor.Err(); err != nil {
        panic(err)

        return
    }

    if err = cursor.Close(context.TODO()); err != nil {
        panic(err)
    }

    log.Printf("Successfully loaded %d servers", len(serversMap))
}
