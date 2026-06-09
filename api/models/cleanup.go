package models

import (
	"context"
	"database/sql"
	"fmt"
	"garden-planner/api/config"
)

type CleanupItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type RenameRequest struct {
	Field string `json:"field"` // "legumes"|"varietes"|"lieux"|"tags"
	From  string `json:"from"`
	To    string `json:"to"`
}

type DeleteRequest struct {
	Field  string `json:"field"`
	Value  string `json:"value"`
	Action string `json:"action"` // "clear" | "replace"
	With   string `json:"with"`   // valeur cible si action="replace"
}

const countExpr = "CASE WHEN COUNT(*) > 20 THEN 21 ELSE COUNT(*) END"

func GetCleanupList(ctx context.Context, gardenId string, field string, legumeFilter string) ([]CleanupItem, error) {
	var rows *sql.Rows
	var err error

	switch field {
	case "legumes":
		rows, err = config.DB.QueryContext(ctx, fmt.Sprintf(`
			SELECT legume, %s FROM action_log
			WHERE jardin_id = ? AND legume != ''
			GROUP BY legume ORDER BY legume ASC`, countExpr), gardenId)
	case "varietes":
		if legumeFilter != "" {
			rows, err = config.DB.QueryContext(ctx, fmt.Sprintf(`
				SELECT variete, %s FROM action_log
				WHERE jardin_id = ? AND variete != '' AND legume = ?
				GROUP BY variete ORDER BY variete ASC`, countExpr), gardenId, legumeFilter)
		} else {
			rows, err = config.DB.QueryContext(ctx, fmt.Sprintf(`
				SELECT variete, %s FROM action_log
				WHERE jardin_id = ? AND variete != ''
				GROUP BY variete ORDER BY variete ASC`, countExpr), gardenId)
		}
	case "lieux":
		rows, err = config.DB.QueryContext(ctx, fmt.Sprintf(`
			SELECT lieu, %s FROM action_log
			WHERE jardin_id = ? AND lieu != ''
			GROUP BY lieu ORDER BY lieu ASC`, countExpr), gardenId)
	case "tags":
		rows, err = config.DB.QueryContext(ctx, fmt.Sprintf(`
			SELECT j.value, %s FROM action_log, json_each(action_log.tags) j
			WHERE jardin_id = ?
			GROUP BY j.value ORDER BY j.value ASC`, countExpr), gardenId)
	default:
		return nil, fmt.Errorf("unknown field: %s", field)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CleanupItem, 0)
	for rows.Next() {
		var item CleanupItem
		if err := rows.Scan(&item.Value, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func RenameCleanupValue(ctx context.Context, gardenId string, req RenameRequest) (int, error) {
	var result sql.Result
	var err error

	switch req.Field {
	case "legumes", "legume":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET legume = ? WHERE legume = ? AND jardin_id = ?`,
			req.To, req.From, gardenId)
	case "varietes", "variete":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET variete = ? WHERE variete = ? AND jardin_id = ?`,
			req.To, req.From, gardenId)
	case "lieux", "lieu":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET lieu = ? WHERE lieu = ? AND jardin_id = ?`,
			req.To, req.From, gardenId)
	case "tags", "tag":
		result, err = config.DB.ExecContext(ctx, `
			UPDATE action_log
			SET tags = (SELECT json_group_array(CASE WHEN j.value = ? THEN ? ELSE j.value END)
			            FROM json_each(tags) j)
			WHERE jardin_id = ? AND EXISTS (SELECT 1 FROM json_each(tags) j WHERE j.value = ?)`,
			req.From, req.To, gardenId, req.From)
	default:
		return 0, fmt.Errorf("unknown field: %s", req.Field)
	}

	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func DeleteCleanupValue(ctx context.Context, gardenId string, req DeleteRequest) (int, error) {
	if req.Action == "replace" {
		return RenameCleanupValue(ctx, gardenId, RenameRequest{Field: req.Field, From: req.Value, To: req.With})
	}

	var result sql.Result
	var err error

	switch req.Field {
	case "legumes", "legume":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET legume = '' WHERE legume = ? AND jardin_id = ?`,
			req.Value, gardenId)
	case "varietes", "variete":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET variete = '' WHERE variete = ? AND jardin_id = ?`,
			req.Value, gardenId)
	case "lieux", "lieu":
		result, err = config.DB.ExecContext(ctx,
			`UPDATE action_log SET lieu = '' WHERE lieu = ? AND jardin_id = ?`,
			req.Value, gardenId)
	case "tags", "tag":
		result, err = config.DB.ExecContext(ctx, `
			UPDATE action_log
			SET tags = (SELECT json_group_array(j.value) FROM json_each(tags) j WHERE j.value != ?)
			WHERE jardin_id = ? AND EXISTS (SELECT 1 FROM json_each(tags) j WHERE j.value = ?)`,
			req.Value, gardenId, req.Value)
	default:
		return 0, fmt.Errorf("unknown field: %s", req.Field)
	}

	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}
