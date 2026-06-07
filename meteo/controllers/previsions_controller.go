package controllers

import (
	"garden-planner/meteo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPrevisions retourne les prévisions horaires (créneaux de 3h) pour la journée en cours.
// Paramètre : station (code station MétéoFrance, ex: 56243001)
func GetPrevisions(c *gin.Context) {
	station := c.Query("station")
	if station == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètre station requis"})
		return
	}

	info, err := models.GetStationCoords(station)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "coordonnées station: " + err.Error()})
		return
	}

	forecasts, err := models.GetPrevisions(info.Lat, info.Lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "prévisions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecasts)
}
