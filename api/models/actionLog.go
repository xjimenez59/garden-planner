package models

import (
	"context"
	"garden-planner/api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActionLog struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ParentId   primitive.ObjectID `json:"_parentId" bson:"_parentId,omitempty"`
	Jardin     string             `json:"jardin" bson:"jardin"`
	DateAction primitive.DateTime `json:"dateAction" bson:"dateAction"`
	Action     string             `json:"action" bson:"action"`
	Statut     string             `json:"statut" bson:"statut"`
	Lieu       string             `json:"lieu" bson:"lieu"`
	Legume     string             `json:"legume" bson:"legume"`
	Variete    string             `json:"variete" bson:"variete"`
	Qte        int                `json:"qte" bson:"qte"`
	Poids      int                `json:"poids" bson:"poids"`
	Notes      string             `json:"notes" bson:"notes"`
	Photos     []string           `json:"photos" bson:"photos"`
	Tags       []string           `json:"tags" bson:"tags"`
}

// Renvoie un slice avec toutes les actions, dans l'ordre chronologique inverse (plus récentes en premier)
func GetLogs(ctx context.Context) (result []ActionLog, err error) {
	result = make([]ActionLog, 0)
	var data *mongo.Cursor

	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{Key: "dateAction", Value: -1}, {Key: "legume", Value: 1}})

	data, err = config.DB.Collection("actionLog").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if err := data.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ActionLog) Save(ctx context.Context) (id primitive.ObjectID, err error) {
	logsCollection := config.DB.Collection("actionLog")
	id = primitive.NilObjectID
	if a.ID.IsZero() {
		var result *mongo.InsertOneResult
		result, err = logsCollection.InsertOne(ctx, a)
		if err == nil {
			id = result.InsertedID.(primitive.ObjectID)
		}
	} else {
		filter := bson.D{{"_id", a.ID}}
		_, err = logsCollection.ReplaceOne(ctx, filter, a)
		if err == nil {
			id = a.ID
		}
	}
	return id, err
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

func DeleteLog(ctx context.Context, id primitive.ObjectID) (err error) {

	logsCollection := config.DB.Collection("actionLog")
	filter := bson.D{{"_id", id}}
	_, err = logsCollection.DeleteOne(ctx, filter)
	return err
}

func GetTags(ctx context.Context) (result []string, err error) {
	result = make([]string, 0)
	var data []interface{}

	filter := bson.D{}
	data, err = config.DB.Collection("actionLog").Distinct(ctx, "tags", filter)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		tag, ok := v.(string)
		if ok {
			result = append(result, tag)
		}
	}
	return result, nil
}
