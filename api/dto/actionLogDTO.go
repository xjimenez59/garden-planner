// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    actionLog, err := UnmarshalActionLog(bytes)
//    bytes, err = actionLog.Marshal()

package dto

import (
	"encoding/json"
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
	Lot        string   `json:"Lot"`
	Legume     string   `json:"legume"`
	Variete    string   `json:"variete"`
	Qte        int      `json:"qte"`
	Poids      int      `json:"poids"`
	Notes      string   `json:"notes"`
	Photos     []string `json:"photos"`
	Tags       []string `json:"tags"`
}
