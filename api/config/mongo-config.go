package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDbuRI = "mongodb://root:zorglub@localhost:27017/?authSource=admin&readPreference=primary&ssl=false&directConnection=true"
var mongoDbName = "GardenPlannerDB"

var DB *mongo.Database = ConnectDatabase()

func ConnectDatabase() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDbuRI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(mongoDbName)
	return db
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(mongoDbName).Collection(collectionName)
	return collection
}
