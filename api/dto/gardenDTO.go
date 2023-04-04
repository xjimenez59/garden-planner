package dto

import (
	"encoding/json"
	"garden-planner/api/models"
)

type GardenDTO struct {
	ID    string `json:"_id"`
	Nom   string `json:"nom"`
	Notes string `json:"notes"`
}

func UnmarshalGardenDTO(data []byte) (GardenDTO, error) {
	var r GardenDTO
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GardenDTO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (d *GardenDTO) FromGardenModel(g models.Garden) {
	d.ID = g.ID.String()
	d.Nom = g.Nom
	d.Notes = g.Notes
}
