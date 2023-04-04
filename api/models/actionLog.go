package models

import (
	"context"
	"garden-planner/api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActionLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	ParentId   primitive.ObjectID `bson:"_parentId"`
	Jardin     string             `bson:"jardin"`
	DateAction primitive.DateTime `bson:"dateAction"`
	Action     string             `bson:"action"`
	Statut     string             `bson:"statut"`
	Lieu       string             `bson:"lieu"`
	Lot        string             `bson:"Lot"`
	Legume     string             `bson:"legume"`
	Variete    string             `bson:"variete"`
	Qte        int                `bson:"qte"`
	Poids      int                `bson:"poids"`
	Notes      string             `bson:"notes"`
	Photos     []string           `bson:"photos"`
	Tags       []string           `bson:"tags"`
}

func GetLogs(ctx context.Context) (result []ActionLog, err error) {
	result = make([]ActionLog, 0)
	var data *mongo.Cursor

	data, err = config.DB.Collection("actionLog").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err := data.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ActionLog) Save(ctx context.Context) (err error) {

	logsCollection := config.DB.Collection("actionLog")
	if a.ID.IsZero() {
		_, err = logsCollection.InsertOne(ctx, a)
	} else {
		_, err = logsCollection.UpdateByID(ctx, a.ID, a)
	}
	return err
}

func SaveLogs(ctx context.Context, logs []ActionLog) (updsertedLogsCount int, err error) {

	logsCollection := config.DB.Collection("actionLog")
	models := []mongo.WriteModel{}

	for _, a := range logs {
		m := mongo.NewReplaceOneModel().
			SetFilter(bson.D{{Key: "_id", Value: a.ID}}).
			SetReplacement(a).
			SetUpsert(true)
		models = append(models, m)
	}

	results, err := logsCollection.BulkWrite(ctx, models)

	return int(results.UpsertedCount), err
}
