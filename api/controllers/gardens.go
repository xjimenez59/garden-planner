package controllers

import (
	"context"
	"garden-planner/api/dto"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetGardens(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var gardens, err = models.GetGardens(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	var result []dto.GardenDTO
	for _, v := range gardens {
		var gdto dto.GardenDTO
		gdto.FromGardenModel(v)
		result = append(result, gdto)
	}

	c.IndentedJSON(http.StatusOK, result)
}
