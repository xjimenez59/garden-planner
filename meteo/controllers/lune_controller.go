package controllers

import (
	"garden-planner/meteo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

	c.JSON(http.StatusOK, gin.H{
		"date":                   date.Format("2006-01-02"),
		"revolution_periodique":  models.RevolutionPeriodique(date),
		"revolution_cyclique":    models.RevolutionCyclique(date),
	})
}
