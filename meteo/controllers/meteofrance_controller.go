package controllers

import (
	"context"
	"garden-planner/meteo/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MF_GetMeteo retourne les données météo journalières depuis la table meteofrance_quotidien.
//
// Paramètres de requête :
//   - station  : identifiant de la station (ex. 56243001)
//   - date_deb : date de début au format YYYYMMDD
//   - date_fin : date de fin au format YYYYMMDD
func MF_GetMeteo(c *gin.Context) {
	station := c.Query("station")
	dateDeb := c.Query("date_deb")
	dateFin := c.Query("date_fin")

	if station == "" || dateDeb == "" || dateFin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètres station, date_deb et date_fin requis"})
		return
	}

	rows, err := models.GetMeteoQuotidien(context.Background(), station, dateDeb, dateFin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// MF_GetMeteoForDates retourne les données météo pour une liste précise de dates.
// Body JSON : {"station": "56243001", "dates": ["20250816", "20240305"]}
// Les dates sont au format YYYYMMDD (identique à la colonne 'date' de la table).
func MF_GetMeteoForDates(c *gin.Context) {
	var body struct {
		Station string   `json:"station"`
		Dates   []string `json:"dates"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body JSON invalide"})
		return
	}
	if body.Station == "" || len(body.Dates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "station et dates requis"})
		return
	}
	// Sanitize : s'assurer que les dates ne contiennent que des chiffres
	for _, d := range body.Dates {
		if strings.ContainsAny(d, " '\";-") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format de date invalide"})
			return
		}
	}
	rows, err := models.GetMeteoQuotidienForDates(context.Background(), body.Station, body.Dates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// MF_CommandeQuotidienne appelle le point d'entrée MétéoFrance
// GET /commande-station/quotidienne et retourne la réponse brute (numéro de commande).
//
// Paramètres de requête :
//   - station   : identifiant de la station (ex. 56243001)
//   - date_deb  : date de début au format YYYYMMDD
//   - date_fin  : date de fin au format YYYYMMDD
func MF_CommandeQuotidienne(c *gin.Context) {
	station := c.Query("station")
	dateDeb := c.Query("date_deb")
	dateFin := c.Query("date_fin")

	if station == "" || dateDeb == "" || dateFin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètres station, date_deb et date_fin requis"})
		return
	}

	body, statusCode, err := models.MFCommandeQuotidienne(station, dateDeb, dateFin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(statusCode, "application/json", body)
}

// MF_GetResultats appelle le point d'entrée MétéoFrance
// GET /commande/fichier, parse le CSV retourné et sauvegarde les données dans SQLite.
//
// Paramètre de requête :
//   - id_cmde : identifiant de commande retourné par MF_CommandeQuotidienne
func MF_GetResultats(c *gin.Context) {
	idCmde := c.Query("id_cmde")
	if idCmde == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramètre id_cmde requis"})
		return
	}

	rows, statusCode, err := models.MFGetFichier(idCmde)
	if err != nil {
		if statusCode == http.StatusAccepted {
			c.JSON(http.StatusAccepted, gin.H{"message": "fichier en cours de préparation, réessayez dans quelques instants"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	for _, row := range rows {
		if err := row.Save(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur sauvegarde: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"saved": len(rows),
	})
}
