// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    actionLog, err := UnmarshalActionLog(bytes)
//    bytes, err = actionLog.Marshal()

package main

import "encoding/json"

func UnmarshalActionLog(data []byte) (ActionLog, error) {
	var r ActionLog
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ActionLog) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ActionLog struct {
	ID         string   `json:"id"`
	Jardin     string   `json:"jardin"`
	DateAction string   `json:"dateAction"`
	Action     string   `json:"action"`
	Statut     string   `json:"statut"`
	Lieu       string   `json:"lieu"`
	Lot        string   `json:"Lot"`
	Legume     string   `json:"legume"`
	Variete    string   `json:"variete"`
	Qte        float64  `json:"qte"`
	Poids      float64  `json:"poids"`
	Notes      string   `json:"notes"`
	Photos     []string `json:"photos"`
}

var actionLogDemoData = []ActionLog{
	{ID: "1", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Patates douces radiateur", Legume: "Patate douce", Variete: "", Qte: 13, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "2", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Poivrons surprise", Legume: "Poivron", Variete: "", Qte: 7, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "3", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Aubergines 2022", Legume: "Aubergine", Variete: "", Qte: 4, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "4", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "", Legume: "Laitue", Variete: "", Qte: 13, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "5", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Fèves 2022", Legume: "Fève", Variete: "", Qte: 6, Poids: 0, Notes: "", Photos: []string{}},
	{ID: "6", Jardin: "Potager Maison", DateAction: "2023-03-13", Action: "Semis", Statut: "Todo", Lieu: "", Lot: "Fèves Ste Marthe 2023", Legume: "Fève", Variete: "", Qte: 15, Poids: 0, Notes: "", Photos: []string{}},
}
