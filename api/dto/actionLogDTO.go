package dto

import (
	"garden-planner/api/models"
)

type ActionLogDTO struct {
	ID         string   `json:"_id"`
	ParentId   string   `json:"_parentId"`
	JardinId   string   `json:"jardinId"`
	DateAction string   `json:"dateAction"`
	Action     string   `json:"action"`
	Statut     string   `json:"statut"`
	Lieu       string   `json:"lieu"`
	Legume     string   `json:"legume"`
	Variete    string   `json:"variete"`
	Qte        int      `json:"qte"`
	Poids      int      `json:"poids"`
	Notes      string   `json:"notes"`
	Photos     []string `json:"photos"`
	Tags       []string `json:"tags"`
}

func (d *ActionLogDTO) FromActionLogModel(a models.ActionLog) {
	d.ID = a.ID
	d.ParentId = a.ParentId
	d.JardinId = a.JardinId
	d.DateAction = a.DateAction
	d.Action = a.Action
	d.Statut = a.Statut
	d.Lieu = a.Lieu
	d.Legume = a.Legume
	d.Variete = a.Variete
	d.Qte = a.Qte
	d.Poids = a.Poids
	d.Notes = a.Notes
	d.Photos = a.Photos
	d.Tags = a.Tags
}

func (d *ActionLogDTO) ToActionLog() (models.ActionLog, error) {
	return models.ActionLog{
		ID:         d.ID,
		ParentId:   d.ParentId,
		JardinId:   d.JardinId,
		DateAction: d.DateAction,
		Action:     d.Action,
		Statut:     d.Statut,
		Lieu:       d.Lieu,
		Legume:     d.Legume,
		Variete:    d.Variete,
		Qte:        d.Qte,
		Poids:      d.Poids,
		Notes:      d.Notes,
		Photos:     d.Photos,
		Tags:       d.Tags,
	}, nil
}
