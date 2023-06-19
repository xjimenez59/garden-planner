package models

import (
	"context"
	"garden-planner/api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Recolte struct {
	Legume string
	Poids  int
	Qte    int
}

// extrait le referentiel des legumes à partir des "actionLog" saisis
func GetRecoltes(ctx context.Context) (result []Recolte, err error) {
	result = make([]Recolte, 0)

	matchStage := bson.D{{"$match", bson.D{{"action", "Récolte"}}}}

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
				{"poids", bson.D{{"$sum", "$poids"}}},
				{"qte", bson.D{{"$sum", "$qte"}}},
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

	cursor, err := config.DB.Collection("actionLog").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage})
	if err != nil {
		return nil, err
	}

	var data []struct {
		//		Famille string `bson: "famille"`
		Legume string `bson: "legume"`
		Poids  int    `bson: "poids"`
		Qte    int    `bson: "qte"`
	}

	if cursor.All(ctx, &data) != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, Recolte{
			Legume: v.Legume,
			Poids:  v.Poids,
			Qte:    v.Qte,
		})
	}
	return result, nil
}
