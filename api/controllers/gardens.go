package controllers

import (
	"context"
	"garden-planner/api/dto"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetGardens(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var gardens, err = models.GetGardens(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	var result = make([]dto.GardenDTO, 0)
	for _, v := range gardens {
		var gdto dto.GardenDTO
		gdto.FromGardenModel(v)
		result = append(result, gdto)
	}
	c.IndentedJSON(http.StatusOK, result)
}

func PostGarden(c *gin.Context) {
	var postedDTO dto.GardenDTO
	var garden models.Garden
	var err error
	var id primitive.ObjectID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	garden, err = postedDTO.ToGardenModel()
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	id, err = garden.Save(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, map[string]string{"_id": id.Hex()})
}

func DeleteGarden(c *gin.Context) {
	var err error
	var id string = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	models.GetGarden(ctx, objId)

	//--- supprimer d'abord tous les logs associés

	err = models.DeleteGarden(ctx, objId)
	if err != nil {
		c.IndentedJSON(http.StatusNotModified, err)
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]string{"_id": id})
}
