package models

import (
	"context"
	"garden-planner/api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecolteAnnee struct {
	Annee int `bson: "annee"`
	Poids int `bson: "poids"`
	Qte   int `bson: "qte"`
}

type Recolte struct {
	Legume string
	Annees []RecolteAnnee `bson: "annees"`
}

// extrait le referentiel des legumes à partir des "actionLog" saisis
func GetRecoltes(ctx context.Context) (result []Recolte, err error) {
	const moisMaxRecolte = 3
	result = make([]Recolte, 0)

	matchStage := bson.D{{"$match", bson.M{"action": "Récolte"}}}

	year02Stage := bson.D{{"$addFields", bson.M{
		"annee": bson.D{{"$toInt", bson.D{{"$substrBytes", bson.A{"$dateAction", 0, 4}}}}},
		"mois":  bson.D{{"$toInt", bson.D{{"$substrBytes", bson.A{"$dateAction", 5, 2}}}}},
	}}}
	year03Stage := bson.D{{"$set", bson.M{"anneeRecolte": bson.D{{"$add", []interface{}{
		"$annee",
		bson.D{{"$cond", []interface{}{
			bson.D{{"$lte", []interface{}{"$mois", moisMaxRecolte}}},
			-1,
			0},
		}},
	},
	}}}}}

	groupStageAnnee := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"legume", "$legume"},
						{"annee", "$anneeRecolte"},
					},
				},
				{"legume", bson.D{{"$first", "$legume"}}},
				{"annee", bson.D{{"$first", "$anneeRecolte"}}},
				{"poids", bson.D{{"$sum", "$poids"}}},
				{"qte", bson.D{{"$sum", "$qte"}}},
			},
		},
	}

	groupStageLegume := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"legume", "$legume"},
					},
				},
				{"legume", bson.D{{"$first", "$legume"}}},
				{"annees", bson.D{{"$addToSet", bson.D{{"annee", "$annee"}, {"poids", "$poids"}, {"qte", "$qte"}}}}},
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

	cursor, err := config.DB.Collection("actionLog").Aggregate(ctx, mongo.Pipeline{
		matchStage,
		year02Stage, year03Stage,
		groupStageAnnee,
		groupStageLegume,
		sortStage})
	if err != nil {
		return nil, err
	}

	var data []struct {
		Legume string         `bson: "legume"`
		Annees []RecolteAnnee `bson: "annees"`
	}

	if cursor.All(ctx, &data) != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, Recolte{
			Legume: v.Legume,
			Annees: v.Annees,
		})
	}
	return result, nil
}
