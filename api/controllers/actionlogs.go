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

	actionLog, err = dtoToActionLog(postedDTO)
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
		a, err = dtoToActionLog(v)
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

func dtoToActionLog(d dto.ActionLogDTO) (a models.ActionLog, err error) {
	a = models.ActionLog{}
	if d.ID == "" {
		a.ID = primitive.NilObjectID
	} else {
		a.ID, err = primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return a, err
		}
	}

	if d.ParentId == "" {
		a.ParentId = primitive.NilObjectID
	} else {
		a.ParentId, err = primitive.ObjectIDFromHex(d.ParentId)
		if err != nil {
			return a, err
		}
	}
	var dateAction time.Time
	dateAction, err = time.Parse("2006-01-02", d.DateAction)
	if err != nil {
		return a, err
	}
	a.Jardin = d.Jardin
	a.DateAction = primitive.NewDateTimeFromTime(dateAction)
	a.Action = d.Action
	a.Statut = d.Statut
	a.Lieu = d.Lieu
	a.Legume = d.Legume
	a.Variete = d.Variete
	a.Qte = d.Qte
	a.Poids = d.Poids
	a.Notes = d.Notes
	a.Photos = d.Photos
	a.Tags = d.Tags

	return a, nil
}
