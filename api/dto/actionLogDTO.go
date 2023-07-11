// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    actionLog, err := UnmarshalActionLog(bytes)
//    bytes, err = actionLog.Marshal()

package dto

import (
	"garden-planner/api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (d *ActionLogDTO) FromGardenModel(a models.ActionLog) {
	d.ID = a.ID.String()
	d.ParentId = a.ParentId.String()
	d.JardinId = a.JardinId.String()
	d.Action = a.Action
	d.DateAction = a.DateAction.Time().Format("2006-01-02")
	d.Legume = a.Legume
	d.Lieu = a.Lieu
	d.Notes = a.Notes
	d.Photos = a.Photos
	d.Poids = a.Poids
	d.Qte = a.Qte
	d.Statut = a.Statut
	d.Tags = a.Tags
	d.Variete = a.Variete

}

func (d *ActionLogDTO) ToActionLog() (a models.ActionLog, err error) {
	a = models.ActionLog{}
	if d.ID == "" {
		a.ID = primitive.NilObjectID
	} else {
		a.ID, err = primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return a, err
		}
	}

	if d.ParentId == "" {
		a.ParentId = primitive.NilObjectID
	} else {
		a.ParentId, err = primitive.ObjectIDFromHex(d.ParentId)
		if err != nil {
			return a, err
		}
	}

	if d.JardinId == "" {
		a.JardinId = primitive.NilObjectID
	} else {
		a.JardinId, err = primitive.ObjectIDFromHex(d.JardinId)
		if err != nil {
			return a, err
		}
	}

	var dateAction time.Time
	dateAction, err = time.Parse("2006-01-02", d.DateAction)
	if err != nil {
		return a, err
	}
	a.DateAction = primitive.NewDateTimeFromTime(dateAction)
	a.Action = d.Action
	a.Statut = d.Statut
	a.Lieu = d.Lieu
	a.Legume = d.Legume
	a.Variete = d.Variete
	a.Qte = d.Qte
	a.Poids = d.Poids
	a.Notes = d.Notes
	a.Photos = d.Photos
	a.Tags = d.Tags

	return a, nil
}
