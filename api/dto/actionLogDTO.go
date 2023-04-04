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

var ActionLogDemoData = []ActionLogDTO{
	{ID: "1", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Patates douces radiateur", Legume: "Patate douce", Variete: "", Qte: 13, Poids: 0, Notes: "", Photos: []string{}, Tags: []string{}},
	{ID: "2", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Poivrons surprise", Legume: "Poivron", Variete: "", Qte: 7, Poids: 0, Notes: "", Photos: []string{}, Tags: []string{}},
	{ID: "3", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Aubergines 2022", Legume: "Aubergine", Variete: "", Qte: 4, Poids: 0, Notes: "", Photos: []string{}, Tags: []string{}},
	{ID: "4", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "", Legume: "Laitue", Variete: "", Qte: 13, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "5", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Fèves 2022", Legume: "Fève", Variete: "", Qte: 6, Poids: 0, Notes: "", Photos: []string{}, Tags: []string{}},
	{ID: "6", ParentId: "", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Fèves Ste Marthe 2023", Legume: "Fève", Variete: "", Qte: 15, Poids: 0, Notes: "", Photos: []string{}, Tags: []string{}},
}
