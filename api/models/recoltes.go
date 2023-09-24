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
	Lieu   string
	Annees []RecolteAnnee `bson: "annees"`
}

// --- Calcule l'année réelle et le mois de récolte, puis l'année d'affectation de la récolte
func AppendYearStages(src mongo.Pipeline, moisMaxRecolte int) (result mongo.Pipeline) {

	result = append(src, bson.D{{"$addFields", bson.M{
		"annee": bson.D{{"$toInt", bson.D{{"$substrBytes", bson.A{"$dateAction", 0, 4}}}}},
		"mois":  bson.D{{"$toInt", bson.D{{"$substrBytes", bson.A{"$dateAction", 5, 2}}}}},
	}}})

	result = append(result, bson.D{{"$set", bson.M{"anneeRecolte": bson.D{{"$add", []interface{}{
		"$annee",
		bson.D{{"$cond", []interface{}{
			bson.D{{"$lte", []interface{}{"$mois", moisMaxRecolte}}},
			-1,
			0},
		}},
	},
	}}}}})

	return result
}

func GetRecoltes(ctx context.Context) (result []Recolte, err error) {
	const moisMaxRecolte = 3
	result = make([]Recolte, 0)

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", bson.M{"action": "Récolte"}}}) //-- filtre les récoltes
	pipe = AppendYearStages(pipe, moisMaxRecolte)                        //-- calcule l'année de récolte

	pipe = append(pipe, bson.D{ //-- cumul par légume et par année
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
	})

	pipe = append(pipe, bson.D{ // regroupe par légume avec une liste d'années
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
	})

	pipe = append(pipe, bson.D{ // ordonne par nom de légume
		{
			"$sort", bson.D{
				{"legume", 1},
			},
		},
	})

	cursor, err := config.DB.Collection("actionLog").Aggregate(ctx, pipe)
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
			Lieu:   "",
			Annees: v.Annees,
		})
	}
	return result, nil
}

func GetRecoltesLieux(ctx context.Context) (result []Recolte, err error) {
	const moisMaxRecolte = 3
	result = make([]Recolte, 0)

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", bson.M{"action": "Récolte"}}}) //-- filtre les récoltes
	pipe = AppendYearStages(pipe, moisMaxRecolte)                        //-- calcule l'année de récolte

	pipe = append(pipe, bson.D{ //-- cumul par lieu, légume et par année
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"lieu", "$lieu"},
						{"legume", "$legume"},
						{"annee", "$anneeRecolte"},
					},
				},
				{"lieu", bson.D{{"$first", "$lieu"}}},
				{"legume", bson.D{{"$first", "$legume"}}},
				{"annee", bson.D{{"$first", "$anneeRecolte"}}},
				{"poids", bson.D{{"$sum", "$poids"}}},
				{"qte", bson.D{{"$sum", "$qte"}}},
			},
		},
	})

	pipe = append(pipe, bson.D{ // regroupe par lieu et légume avec une liste d'années
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"lieu", "$lieu"},
						{"legume", "$legume"},
					},
				},
				{"lieu", bson.D{{"$first", "$lieu"}}},
				{"legume", bson.D{{"$first", "$legume"}}},
				{"annees", bson.D{{"$addToSet", bson.D{{"annee", "$annee"}, {"poids", "$poids"}, {"qte", "$qte"}}}}},
			},
		},
	})

	pipe = append(pipe, bson.D{ // ordonne par nom de légume
		{
			"$sort", bson.D{
				{"lieu", 1},
				{"legume", 1},
			},
		},
	})

	cursor, err := config.DB.Collection("actionLog").Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var data []struct {
		Lieu   string         `bson: "lieu"`
		Legume string         `bson: "legume"`
		Annees []RecolteAnnee `bson: "annees"`
	}

	if cursor.All(ctx, &data) != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, Recolte{
			Legume: v.Legume,
			Lieu:   v.Lieu,
			Annees: v.Annees,
		})
	}
	return result, nil
}
