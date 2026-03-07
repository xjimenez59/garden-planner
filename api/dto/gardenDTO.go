package dto

import (
	"garden-planner/api/models"
)

type GardenRoleDTO struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

type GardenDTO struct {
	ID             string          `json:"_id"`
	Nom            string          `json:"nom"`
	Notes          string          `json:"notes"`
	MoisFinRecolte int             `json:"moisFinRecolte"`
	MoisFinSemis   int             `json:"moisFinSemis"`
	Localisation   string          `json:"localisation"`
	Surface        int             `json:"surface"`
	MeteofSite     string          `json:"meteofSite"`
	Jardiniers     []GardenRoleDTO `json:"jardiniers"`
}

func (d *GardenDTO) FromGardenModel(g models.Garden) {
	d.ID = g.ID
	d.Nom = g.Nom
	d.Notes = g.Notes
	d.MoisFinRecolte = g.MoisFinRecolte
	d.MoisFinSemis = g.MoisFinSemis
	d.Localisation = g.Localisation
	d.Surface = g.Surface
	d.MeteofSite = g.MeteofSite
	d.Jardiniers = make([]GardenRoleDTO, 0)
	for _, v := range g.Jardiniers {
		d.Jardiniers = append(d.Jardiniers, GardenRoleDTO{UserID: v.UserID, Role: v.Role})
	}
}

func (d *GardenDTO) ToGardenModel() (models.Garden, error) {
	g := models.Garden{
		ID:             d.ID,
		Nom:            d.Nom,
		Notes:          d.Notes,
		MoisFinRecolte: d.MoisFinRecolte,
		MoisFinSemis:   d.MoisFinSemis,
		Localisation:   d.Localisation,
		Surface:        d.Surface,
		MeteofSite:     d.MeteofSite,
		Jardiniers:     make([]models.GardenRole, 0),
	}
	for _, v := range d.Jardiniers {
		g.Jardiniers = append(g.Jardiniers, models.GardenRole{UserID: v.UserID, Role: v.Role})
	}
	return g, nil
}
