package controllers

import (
	"context"
	"garden-planner/meteo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMeteo(c *gin.Context) {

	site := c.Param("site")
	date := c.Param("date")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var previsions, err = models.GetMeteo(ctx, site, date)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, previsions)
}
