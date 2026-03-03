package models

import (
	"context"
	"database/sql"
	"garden-planner/api/config"
	"strings"
)

type Legume struct {
	ID      string   `json:"_id"`
	Nom     string   `json:"nom"`
	Famille string   `json:"famille"`
	Varietes []string `json:"variete"`
	Notes   string   `json:"notes"`
}

// extrait le référentiel des légumes à partir des actionLog saisis
func GetLegumesFromActions(ctx context.Context) ([]Legume, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT legume, GROUP_CONCAT(DISTINCT variete) AS varietes
		FROM action_log
		WHERE legume != ''
		GROUP BY legume
		ORDER BY legume`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legumes := make([]Legume, 0)
	for rows.Next() {
		var l Legume
		var varietesStr sql.NullString
		if err := rows.Scan(&l.Nom, &varietesStr); err != nil {
			return nil, err
		}
		l.Varietes = []string{}
		if varietesStr.Valid && varietesStr.String != "" {
			l.Varietes = strings.Split(varietesStr.String, ",")
		}
		legumes = append(legumes, l)
	}
	return legumes, rows.Err()
}
