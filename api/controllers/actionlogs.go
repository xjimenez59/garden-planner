package controllers

import (
	"context"
	"fmt"
	"garden-planner/api/dto"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetLogs(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var actionLogs, err = models.GetLogs(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, actionLogs)
}

func PostLog(c *gin.Context) {
	var postedDTO dto.ActionLogDTO
	var actionLog models.ActionLog
	var err error
	var id primitive.ObjectID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	actionLog, err = postedDTO.ToActionLog()
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	id, err = actionLog.Save(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, map[string]string{"_id": id.Hex()})
}

func PostLogs(c *gin.Context) {
	var postedDTO []dto.ActionLogDTO
	var logActions []models.ActionLog
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	if err = c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	for _, v := range postedDTO {
		var a models.ActionLog
		a, err = v.ToActionLog()
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err)
			return
		}
		logActions = append(logActions, a)
	}
	var updatedLogs = 0
	updatedLogs, err = models.SaveLogs(ctx, logActions)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("{updated: %d}", updatedLogs))
}

func DeleteLog(c *gin.Context) {
	var err error
	var id string = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	a := models.GetLog(ctx, objId)

	err = models.DeleteLog(ctx, objId)
	if err != nil {
		c.IndentedJSON(http.StatusNotModified, err)
		return
	}

	DeleteBucketObjects(c, a.Photos)
	fmt.Println(a.Photos)
	c.IndentedJSON(http.StatusOK, map[string]string{"_id": id})
}

func GetTags(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tags, err = models.GetTags(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, tags)
}

func GetLieux(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var lieux, err = models.GetLieux(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, lieux)
}
