// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    actionLog, err := UnmarshalActionLog(bytes)
//    bytes, err = actionLog.Marshal()

package dto

import (
	"encoding/json"
	"garden-planner/api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalActionLogDTO(data []byte) (ActionLogDTO, error) {
	var r ActionLogDTO
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ActionLogDTO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ActionLogDTO struct {
	ID         string   `json:"_id"`
	ParentId   string   `json:"_parentId"`
	Jardin     string   `json:"jardin"`
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
	var dateAction time.Time
	dateAction, err = time.Parse("2006-01-02", d.DateAction)
	if err != nil {
		return a, err
	}
	a.Jardin = d.Jardin
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
