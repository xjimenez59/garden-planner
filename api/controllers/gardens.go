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

	userID := c.Request.Header.Get("Authorization")

	gardens, err := models.GetGardens(ctx, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := make([]dto.GardenDTO, 0)
	for _, v := range gardens {
		var gdto dto.GardenDTO
		gdto.FromGardenModel(v)
		result = append(result, gdto)
	}
	c.IndentedJSON(http.StatusOK, result)
}

func PostGarden(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var postedDTO dto.GardenDTO
	if err := c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	garden, err := postedDTO.ToGardenModel()
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	id, err := garden.Save(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, map[string]string{"_id": id})
}

func DeleteGarden(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	if _, err := models.GetGarden(ctx, id); err != nil {
		c.IndentedJSON(http.StatusNotFound, err.Error())
		return
	}

	if err := models.DeleteGarden(ctx, id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]string{"_id": id})
}
