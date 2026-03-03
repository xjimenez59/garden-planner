package models

import (
	"context"
	"garden-planner/api/config"

	"github.com/google/uuid"
)

type GardenRole struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

type Garden struct {
	ID             string       `json:"_id"`
	Nom            string       `json:"nom"`
	Notes          string       `json:"notes"`
	MoisFinRecolte int          `json:"moisFinRecolte"`
	MoisFinSemis   int          `json:"moisFinSemis"`
	Localisation   string       `json:"localisation"`
	Surface        int          `json:"surface"`
	Jardiniers     []GardenRole `json:"jardiniers"`
}

func GetGardens(ctx context.Context, userID string) ([]Garden, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT DISTINCT g.id, g.nom, g.notes, g.mois_fin_recolte, g.mois_fin_semis, g.localisation, g.surface
		FROM garden g
		INNER JOIN garden_jardinier gj ON g.id = gj.garden_id
		WHERE gj.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gardens := make([]Garden, 0)
	for rows.Next() {
		var g Garden
		if err := rows.Scan(&g.ID, &g.Nom, &g.Notes, &g.MoisFinRecolte, &g.MoisFinSemis, &g.Localisation, &g.Surface); err != nil {
			return nil, err
		}
		g.Jardiniers, err = getJardiniers(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		gardens = append(gardens, g)
	}
	return gardens, rows.Err()
}

func getJardiniers(ctx context.Context, gardenID string) ([]GardenRole, error) {
	rows, err := config.DB.QueryContext(ctx, `SELECT user_id, role FROM garden_jardinier WHERE garden_id = ?`, gardenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]GardenRole, 0)
	for rows.Next() {
		var r GardenRole
		if err := rows.Scan(&r.UserID, &r.Role); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, rows.Err()
}

func GetGarden(ctx context.Context, id string) (Garden, error) {
	var g Garden
	err := config.DB.QueryRowContext(ctx, `
		SELECT id, nom, notes, mois_fin_recolte, mois_fin_semis, localisation, surface
		FROM garden WHERE id = ?`, id).
		Scan(&g.ID, &g.Nom, &g.Notes, &g.MoisFinRecolte, &g.MoisFinSemis, &g.Localisation, &g.Surface)
	if err != nil {
		return g, err
	}
	g.Jardiniers, err = getJardiniers(ctx, id)
	return g, err
}

func (g *Garden) Save(ctx context.Context) (string, error) {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}

	tx, err := config.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO garden (id, nom, notes, mois_fin_recolte, mois_fin_semis, localisation, surface)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			nom = excluded.nom,
			notes = excluded.notes,
			mois_fin_recolte = excluded.mois_fin_recolte,
			mois_fin_semis = excluded.mois_fin_semis,
			localisation = excluded.localisation,
			surface = excluded.surface`,
		g.ID, g.Nom, g.Notes, g.MoisFinRecolte, g.MoisFinSemis, g.Localisation, g.Surface)
	if err != nil {
		return "", err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM garden_jardinier WHERE garden_id = ?`, g.ID); err != nil {
		return "", err
	}

	for _, j := range g.Jardiniers {
		if _, err = tx.ExecContext(ctx, `INSERT INTO garden_jardinier (garden_id, user_id, role) VALUES (?, ?, ?)`,
			g.ID, j.UserID, j.Role); err != nil {
			return "", err
		}
	}

	return g.ID, tx.Commit()
}

func DeleteGarden(ctx context.Context, id string) error {
	_, err := config.DB.ExecContext(ctx, `DELETE FROM garden WHERE id = ?`, id)
	return err
}
