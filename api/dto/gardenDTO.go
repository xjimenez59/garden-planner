package dto

import (
	"garden-planner/api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GardenDTO struct {
	ID             string `json:"_id"`
	Nom            string `json:"nom"`
	Notes          string `json:"notes"`
	MoisFinRecolte int    `json:"moisFinRecolte"`
	MoisFinSemis   int    `json:"moisFinSemis"`
	Localisation   string `json:"localisation"`
	Surface        int    `json:"surface"`
}

func (d *GardenDTO) FromGardenModel(g models.Garden) {
	d.ID = g.ID.String()
	d.Nom = g.Nom
	d.Notes = g.Notes
	d.MoisFinRecolte = g.MoisFinRecolte
	d.MoisFinSemis = g.MoisFinSemis
	d.Localisation = g.Localisation
	d.Surface = g.Surface
}

func (d *GardenDTO) ToGardenModel() (g models.Garden, err error) {
	g = models.Garden{}
	if d.ID == "" {
		g.ID = primitive.NilObjectID
	} else {
		g.ID, err = primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return g, err
		}
	}
	g.Nom = d.Nom
	g.Notes = d.Notes
	g.MoisFinRecolte = d.MoisFinRecolte
	g.MoisFinSemis = d.MoisFinSemis
	g.Localisation = d.Localisation
	g.Surface = d.Surface

	return g, nil
}
