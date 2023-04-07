package models

import (
	"context"
	"fmt"
	"garden-planner/api/config"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Legume struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Nom      string             `json:"nom" bson:"nom"`
	Famille  string             `json:"famille" bson:"famille"`
	Varietes []string           `json:"variete" bson:"variete"`
	Notes    string             `json:"notes" bson:"notes"`
}

// Renvoie un slice avec tous les legumes mentionnés dans les actions, dansl'ordre alphabétique
func GetLegumes(ctx context.Context) (result []Legume, err error) {
	result = make([]Legume, 0)

	filter := bson.D{}
	opts := options.Distinct()

	legumes, err := config.DB.Collection("actionLog").Distinct(ctx, "legume", filter, opts)
	if err != nil {
		return nil, err
	}

	for _, v := range legumes {
		result = append(result, Legume{
			ID:       primitive.NilObjectID,
			Famille:  "",
			Nom:      fmt.Sprint(v),
			Varietes: []string{},
			Notes:    ""})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Nom < result[j].Nom
	})

	return result, nil
}

// extrait le referentiel des legumes à partir des "actionLog" saisis
func GetLegumesFromActions(ctx context.Context) (result []Legume, err error) {
	result = make([]Legume, 0)

	groupStage := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"famille", "$famille"},
						{"legume", "$legume"},
						//	{"variete", "$variete"},
					},
				},
				{"famille", bson.D{{"$first", "$famille"}}},
				{"legume", bson.D{{"$first", "$legume"}}},
				{"varietes", bson.D{{"$addToSet", "$variete"}}},
			},
		},
	}

	sortStage := bson.D{
		{
			"$sort", bson.D{
				{"legume", 1},
			},
		},
	}

	cursor, err := config.DB.Collection("actionLog").Aggregate(ctx, mongo.Pipeline{groupStage, sortStage})
	if err != nil {
		return nil, err
	}

	var data []struct {
		//		Famille string `bson: "famille"`
		Legume   string   `bson: "legume"`
		Varietes []string `bson: "variete"`
	}

	if cursor.All(ctx, &data) != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, Legume{
			ID:       primitive.NilObjectID,
			Famille:  "",
			Nom:      v.Legume,
			Varietes: v.Varietes,
			Notes:    ""})
	}
	return result, nil
}
