package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- l'IP est remplaée par le nom du conteneur : mongo
var mongoHost = os.Getenv("MONGO_HOST")
var mongoPort = os.Getenv("MONGO_PORT")
var mongoUser = os.Getenv("MONGO_USER")
var mongoPwd = os.Getenv("MONGO_PWD")
var mongoDbName = os.Getenv("MONGO_DBNAME") // GardenPlannerDB

var mongoDbuRI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin&readPreference=primary&ssl=false&directConnection=true", mongoUser, mongoPwd, mongoHost, mongoPort)

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
