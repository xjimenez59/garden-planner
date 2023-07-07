package models

import (
	"context"
	"garden-planner/api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

type Garden struct {
	ID             primitive.ObjectID `bson:"_id"`
	Nom            string             `bson:"nom"`
	Notes          string             `bson:"notes"`
	MoisFinRecolte int                `bson:"moisFinRecolte"`
	MoisFinSemis   int                `bson:"moisFinSemis"`
	Localisation   string             `bson:"localisation"`
	Surface        int                `bson:"surface"`
}

func GetGardens(ctx context.Context) (result []Garden, err error) {
	result = make([]Garden, 0)
	var data *mongo.Cursor

	data, err = config.DB.Collection("garden").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err := data.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetGarden(ctx context.Context, id primitive.ObjectID) (result Garden) {

	gardensCollection := config.DB.Collection("garden")
	filter := bson.D{{"_id", id}}
	var found *mongo.SingleResult
	found = gardensCollection.FindOne(ctx, filter)
	found.Decode(&result)
	return result
}

func (g *Garden) Save(ctx context.Context) (id primitive.ObjectID, err error) {
	gardensCollection := config.DB.Collection("garden")
	id = primitive.NilObjectID
	if g.ID.IsZero() {
		g.ID = primitive.NewObjectID()
		var result *mongo.InsertOneResult
		result, err = gardensCollection.InsertOne(ctx, g)
		if err == nil {
			id = result.InsertedID.(primitive.ObjectID)
		}
	} else {
		filter := bson.D{{"_id", g.ID}}
		_, err = gardensCollection.ReplaceOne(ctx, filter, g)
		if err == nil {
			id = g.ID
		}
	}
	return id, err
}

func DeleteGarden(ctx context.Context, id primitive.ObjectID) (err error) {

	gardensCollection := config.DB.Collection("garden")
	filter := bson.D{{"_id", id}}
	_, err = gardensCollection.DeleteOne(ctx, filter)
	return err
}
