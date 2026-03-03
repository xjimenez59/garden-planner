package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"garden-planner/api/config"
	"time"

	"github.com/google/uuid"
)

type ActionLog struct {
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

func GetLogs(ctx context.Context, jardinId string) ([]ActionLog, error) {
	pastDate := time.Now().AddDate(0, -15, 0).Format("2006-01-02")

	var rows *sql.Rows
	var err error

	if jardinId == "" {
		rows, err = config.DB.QueryContext(ctx, `
			SELECT id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags
			FROM action_log
			WHERE date_action >= ?
			ORDER BY date_action DESC, legume ASC`, pastDate)
	} else {
		rows, err = config.DB.QueryContext(ctx, `
			SELECT id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags
			FROM action_log
			WHERE jardin_id = ? AND date_action >= ?
			ORDER BY date_action DESC, legume ASC`, jardinId, pastDate)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanActionLogs(rows)
}

func scanActionLogs(rows *sql.Rows) ([]ActionLog, error) {
	logs := make([]ActionLog, 0)
	for rows.Next() {
		var a ActionLog
		var photosJSON, tagsJSON string
		if err := rows.Scan(&a.ID, &a.ParentId, &a.JardinId, &a.DateAction, &a.Action, &a.Statut, &a.Lieu, &a.Legume, &a.Variete, &a.Qte, &a.Poids, &a.Notes, &photosJSON, &tagsJSON); err != nil {
			return nil, err
		}
		a.Photos = []string{}
		a.Tags = []string{}
		json.Unmarshal([]byte(photosJSON), &a.Photos)
		json.Unmarshal([]byte(tagsJSON), &a.Tags)
		logs = append(logs, a)
	}
	return logs, rows.Err()
}

func GetLog(ctx context.Context, id string) (ActionLog, error) {
	var a ActionLog
	var photosJSON, tagsJSON string
	err := config.DB.QueryRowContext(ctx, `
		SELECT id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags
		FROM action_log WHERE id = ?`, id).
		Scan(&a.ID, &a.ParentId, &a.JardinId, &a.DateAction, &a.Action, &a.Statut, &a.Lieu, &a.Legume, &a.Variete, &a.Qte, &a.Poids, &a.Notes, &photosJSON, &tagsJSON)
	if err != nil {
		return a, err
	}
	a.Photos = []string{}
	a.Tags = []string{}
	json.Unmarshal([]byte(photosJSON), &a.Photos)
	json.Unmarshal([]byte(tagsJSON), &a.Tags)
	return a, nil
}

func (a *ActionLog) Save(ctx context.Context) (string, error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	photosJSON, _ := json.Marshal(a.Photos)
	tagsJSON, _ := json.Marshal(a.Tags)

	_, err := config.DB.ExecContext(ctx, `
		INSERT INTO action_log (id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			parent_id = excluded.parent_id,
			jardin_id = excluded.jardin_id,
			date_action = excluded.date_action,
			action = excluded.action,
			statut = excluded.statut,
			lieu = excluded.lieu,
			legume = excluded.legume,
			variete = excluded.variete,
			qte = excluded.qte,
			poids = excluded.poids,
			notes = excluded.notes,
			photos = excluded.photos,
			tags = excluded.tags`,
		a.ID, a.ParentId, a.JardinId, a.DateAction, a.Action, a.Statut, a.Lieu, a.Legume, a.Variete, a.Qte, a.Poids, a.Notes, string(photosJSON), string(tagsJSON))
	return a.ID, err
}

func SaveLogs(ctx context.Context, logs []ActionLog) (int, error) {
	tx, err := config.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	updated := 0
	for _, a := range logs {
		if a.ID == "" {
			a.ID = uuid.New().String()
		}
		photosJSON, _ := json.Marshal(a.Photos)
		tagsJSON, _ := json.Marshal(a.Tags)

		res, err := tx.ExecContext(ctx, `
			INSERT INTO action_log (id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				parent_id = excluded.parent_id,
				jardin_id = excluded.jardin_id,
				date_action = excluded.date_action,
				action = excluded.action,
				statut = excluded.statut,
				lieu = excluded.lieu,
				legume = excluded.legume,
				variete = excluded.variete,
				qte = excluded.qte,
				poids = excluded.poids,
				notes = excluded.notes,
				photos = excluded.photos,
				tags = excluded.tags`,
			a.ID, a.ParentId, a.JardinId, a.DateAction, a.Action, a.Statut, a.Lieu, a.Legume, a.Variete, a.Qte, a.Poids, a.Notes, string(photosJSON), string(tagsJSON))
		if err != nil {
			return 0, err
		}
		n, _ := res.RowsAffected()
		updated += int(n)
	}
	return updated, tx.Commit()
}

func DeleteLog(ctx context.Context, id string) error {
	_, err := config.DB.ExecContext(ctx, `DELETE FROM action_log WHERE id = ?`, id)
	return err
}

func GetTags(ctx context.Context) ([]string, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT DISTINCT value
		FROM action_log, json_each(action_log.tags)
		ORDER BY value`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func GetLieux(ctx context.Context) ([]string, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT DISTINCT lieu FROM action_log WHERE lieu != '' ORDER BY lieu`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lieux := make([]string, 0)
	for rows.Next() {
		var lieu string
		if err := rows.Scan(&lieu); err != nil {
			return nil, err
		}
		lieux = append(lieux, lieu)
	}
	return lieux, rows.Err()
}

func UpdateLogsSetGarden(ctx context.Context, newValue string) (int, error) {
	result, err := config.DB.ExecContext(ctx, `UPDATE action_log SET jardin_id = ?`, newValue)
	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}
