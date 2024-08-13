package server

import (
    "context"
    "github.com/holypvp/primal/server/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
)

type ServerService struct{}

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
}

func (service *ServerService) LoadServers() {

}
