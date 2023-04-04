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
	ID    primitive.ObjectID `bson:"_id"`
	Nom   string             `bson:"nom"`
	Notes string             `bson:"notes"`
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

/* var gardenDemoData = []Garden{
	{ID: "1", Nom: "Potager Jactez", Notes: ""},
	{ID: "2", Nom: "Jardin Partagé Tropark", Notes: ""},
}
*/
