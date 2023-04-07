package controllers

import (
	"context"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// -- renvoie le référentiel des légumes
// -- Pour l'instant, extrait le référentie des actions en base
func GetLegumes(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var legumes, err = models.GetLegumesFromActions(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, legumes)
}
