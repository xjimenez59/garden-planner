package controllers

import (
	"garden-planner/meteo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLuneForDates retourne le jour biodynamique pour une liste précise de dates.
// Body JSON : {"dates": ["2025-08-16", "2024-03-05"]}
func GetLuneForDates(c *gin.Context) {
	var body struct {
		Dates []string `json:"dates"`
	}
	if err := c.BindJSON(&body); err != nil || len(body.Dates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body JSON invalide : {\"dates\": [\"YYYY-MM-DD\", ...]}"})
		return
	}

	type luneDay struct {
		Date             string `json:"date"`
		JourBiodynamique string `json:"jour_biodynamique"`
		SigneZodiaque    string `json:"signe_zodiaque"`
	}

	result := make([]luneDay, 0, len(body.Dates))
	for _, ds := range body.Dates {
		d, err := time.Parse("2006-01-02", ds)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format invalide, utiliser YYYY-MM-DD : " + ds})
			return
		}
		jourBio, signe := models.JourBiodynamique(d)
		result = append(result, luneDay{Date: ds, JourBiodynamique: jourBio, SigneZodiaque: signe})
	}
	c.JSON(http.StatusOK, result)
}

// GetLuneRange retourne le jour biodynamique pour chaque jour d'une plage de dates.
// Paramètres : date_deb, date_fin (YYYY-MM-DD)
func GetLuneRange(c *gin.Context) {
	debStr := c.Query("date_deb")
	finStr := c.Query("date_fin")
	if debStr == "" || finStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètres date_deb et date_fin requis (YYYY-MM-DD)"})
		return
	}
	deb, err := time.Parse("2006-01-02", debStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format date_deb invalide"})
		return
	}
	fin, err := time.Parse("2006-01-02", finStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format date_fin invalide"})
		return
	}

	type luneDay struct {
		Date             string `json:"date"`
		JourBiodynamique string `json:"jour_biodynamique"`
		SigneZodiaque    string `json:"signe_zodiaque"`
	}

	result := make([]luneDay, 0)
	for d := deb; !d.After(fin); d = d.AddDate(0, 0, 1) {
		jourBio, signe := models.JourBiodynamique(d)
		result = append(result, luneDay{
			Date:             d.Format("2006-01-02"),
			JourBiodynamique: jourBio,
			SigneZodiaque:    signe,
		})
	}
	c.JSON(http.StatusOK, result)
}

// GetLune retourne les informations sur le cycle lunaire pour une date donnée :
//   - revolution_periodique : "lune_montante" ou "lune_descendante" (cycle tropical ~27,32 j)
//   - revolution_cyclique   : "lune_croissante" ou "lune_decroissante" (cycle synodique ~29,53 j)
//
// Paramètre de requête :
//   - date : format YYYY-MM-DD (optionnel, défaut = aujourd'hui)
func GetLune(c *gin.Context) {
	dateStr := c.Query("date")

	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format de date invalide, utiliser YYYY-MM-DD"})
			return
		}
	}

	jourBio, signe := models.JourBiodynamique(date)
	etatOrbite, perigee, apogee := models.ApogeePerigee(date)
	omegaDeg, noeudAsc, noeudDsc := models.NoeudsLunaires(date)

	c.JSON(http.StatusOK, gin.H{
		"date":                      date.Format("2006-01-02"),
		"revolution_periodique":     models.RevolutionPeriodique(date),
		"revolution_cyclique":       models.RevolutionCyclique(date),
		"jour_biodynamique":         jourBio,
		"signe_zodiaque":            signe,
		"noeud_ascendant_longitude": omegaDeg,
		"prochain_noeud_ascendant":  noeudAsc,
		"prochain_noeud_descendant": noeudDsc,
		"prochain_perigee":          perigee,
		"prochain_apogee":           apogee,
		"etat_orbite":               etatOrbite,
	})
}
