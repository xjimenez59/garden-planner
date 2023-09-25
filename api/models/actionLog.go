package models

import (
	"context"
	"garden-planner/api/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActionLog struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ParentId   primitive.ObjectID `json:"_parentId" bson:"_parentId,omitempty"`
	JardinId   primitive.ObjectID `json:"jardinId" bson:"jardinId"`
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

func GetLogs(ctx context.Context, gardenId primitive.ObjectID) (result []ActionLog, err error) {
	filter := bson.M{}
	if !gardenId.IsZero() {
		filter["jardinId"] = gardenId
	}

	pastDate := primitive.NewDateTimeFromTime(time.Now().AddDate(0, -15, 0))
	filter["dateAction"] = bson.M{"$gte": pastDate}

	return GetLogsFiltered(ctx, filter)
}

// Renvoie un slice avec toutes les actions, dans l'ordre chronologique inverse (plus récentes en premier)
func GetLogsFiltered(ctx context.Context, filter bson.M) (result []ActionLog, err error) {
	result = make([]ActionLog, 0)
	var data *mongo.Cursor

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

func GetLog(ctx context.Context, id primitive.ObjectID) (result ActionLog) {

	logsCollection := config.DB.Collection("actionLog")
	filter := bson.D{{"_id", id}}
	var found *mongo.SingleResult
	found = logsCollection.FindOne(ctx, filter)
	found.Decode(&result)
	return result
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
		if a.ID.IsZero() {
			m := mongo.NewInsertOneModel().SetDocument(a)
			models = append(models, m)
		} else {
			m := mongo.NewReplaceOneModel().
				SetFilter(bson.D{{"_id", a.ID}}).
				SetReplacement(a).
				SetUpsert(false)
			models = append(models, m)
		}
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

func GetLieux(ctx context.Context) (result []string, err error) {
	result = make([]string, 0)
	var data []interface{}

	filter := bson.D{}
	data, err = config.DB.Collection("actionLog").Distinct(ctx, "lieu", filter)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		lieu, ok := v.(string)
		if ok {
			result = append(result, lieu)
		}
	}
	return result, nil
}

func UpdateLogsSetGarden(ctx context.Context, newValue primitive.ObjectID) (updated int, err error) {
	filter := bson.D{}
	logsCollection := config.DB.Collection("actionLog")
	update := bson.M{"$set": bson.M{"jardinId": newValue}, "$unset": bson.M{"jardin": ""}}

	var result *mongo.UpdateResult
	result, err = logsCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return int(result.ModifiedCount), nil

}
