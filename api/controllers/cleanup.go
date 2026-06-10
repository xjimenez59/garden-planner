package controllers

import (
	"context"
	"garden-planner/api/data"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCleanupList(c *gin.Context) {
	gardenId := c.Param("gardenId")
	field := c.Param("field")
	legumeFilter := c.Query("legume")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	items, err := models.GetCleanupList(ctx, gardenId, field, legumeFilter)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, items)
}

func RenameCleanupValue(c *gin.Context) {
	gardenId := c.Param("gardenId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.RenameRequest
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	n, err := models.RenameCleanupValue(ctx, gardenId, req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"updated": n})
}

func DeleteCleanupValue(c *gin.Context) {
	gardenId := c.Param("gardenId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.DeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	n, err := models.DeleteCleanupValue(ctx, gardenId, req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"updated": n})
}

func GetLegumesReference(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", data.LegumesReference)
}
