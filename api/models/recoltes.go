package models

import (
	"context"
	"garden-planner/api/config"
)

type RecolteAnnee struct {
	Annee int
	Poids int
	Qte   int
}

type Recolte struct {
	Legume string
	Lieu   string
	Annees []RecolteAnnee
}

// GetRecoltes retourne les cumuls de récolte par légume et par année.
// Les récoltes de janvier-mars sont rattachées à l'année précédente.
func GetRecoltes(ctx context.Context) ([]Recolte, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT
			legume,
			CASE
				WHEN CAST(strftime('%m', date_action) AS INTEGER) <= 3
				THEN CAST(strftime('%Y', date_action) AS INTEGER) - 1
				ELSE CAST(strftime('%Y', date_action) AS INTEGER)
			END AS annee_recolte,
			SUM(poids) AS poids,
			SUM(qte) AS qte
		FROM action_log
		WHERE action = 'Récolte' AND legume != ''
		GROUP BY legume, annee_recolte
		ORDER BY legume, annee_recolte`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recolteMap := make(map[string]*Recolte)
	order := make([]string, 0)

	for rows.Next() {
		var legume string
		var annee, poids, qte int
		if err := rows.Scan(&legume, &annee, &poids, &qte); err != nil {
			return nil, err
		}
		if _, ok := recolteMap[legume]; !ok {
			recolteMap[legume] = &Recolte{Legume: legume, Annees: []RecolteAnnee{}}
			order = append(order, legume)
		}
		recolteMap[legume].Annees = append(recolteMap[legume].Annees, RecolteAnnee{Annee: annee, Poids: poids, Qte: qte})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]Recolte, 0, len(order))
	for _, k := range order {
		result = append(result, *recolteMap[k])
	}
	return result, nil
}

// GetRecoltesLieux retourne les cumuls de récolte par lieu, légume et par année.
func GetRecoltesLieux(ctx context.Context) ([]Recolte, error) {
	rows, err := config.DB.QueryContext(ctx, `
		SELECT
			lieu,
			legume,
			CASE
				WHEN CAST(strftime('%m', date_action) AS INTEGER) <= 3
				THEN CAST(strftime('%Y', date_action) AS INTEGER) - 1
				ELSE CAST(strftime('%Y', date_action) AS INTEGER)
			END AS annee_recolte,
			SUM(poids) AS poids,
			SUM(qte) AS qte
		FROM action_log
		WHERE action = 'Récolte' AND legume != ''
		GROUP BY lieu, legume, annee_recolte
		ORDER BY lieu, legume, annee_recolte`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type key struct{ lieu, legume string }
	recolteMap := make(map[key]*Recolte)
	order := make([]key, 0)

	for rows.Next() {
		var lieu, legume string
		var annee, poids, qte int
		if err := rows.Scan(&lieu, &legume, &annee, &poids, &qte); err != nil {
			return nil, err
		}
		k := key{lieu, legume}
		if _, ok := recolteMap[k]; !ok {
			recolteMap[k] = &Recolte{Legume: legume, Lieu: lieu, Annees: []RecolteAnnee{}}
			order = append(order, k)
		}
		recolteMap[k].Annees = append(recolteMap[k].Annees, RecolteAnnee{Annee: annee, Poids: poids, Qte: qte})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]Recolte, 0, len(order))
	for _, k := range order {
		result = append(result, *recolteMap[k])
	}
	return result, nil
}
