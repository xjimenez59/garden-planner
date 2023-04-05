package models

import (
	"context"
	"fmt"
	"garden-planner/api/config"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Legume struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Nom   string             `json:"nom" bson:"nom"`
	Notes string             `json:"notes" bson:"notes"`
}

// Renvoie un slice avec toutes les legumes mentionnés dans les actions, dansl'ordre alphabétique
func GetLegumes(ctx context.Context) (result []Legume, err error) {
	result = make([]Legume, 0)

	filter := bson.D{}
	opts := options.Distinct()

	legumes, err := config.DB.Collection("actionLog").Distinct(ctx, "legume", filter, opts)
	if err != nil {
		return nil, err
	}

	for _, v := range legumes {
		result = append(result, Legume{ID: primitive.NilObjectID, Nom: fmt.Sprint(v), Notes: ""})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Nom < result[j].Nom
	})

	return result, nil
}
