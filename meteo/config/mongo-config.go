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

	opts := options.Client().ApplyURI(mongoDbuRI)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal((err))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(mongoDbName)
	return db
}

func CloseDatabase() {
	err := DB.Client().Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

}
