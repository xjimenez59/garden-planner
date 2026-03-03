package controllers

import (
	"context"
	"fmt"
	"garden-planner/api/dto"
	"garden-planner/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetLogs(c *gin.Context) {
	jardinId := c.Param("gardenId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	actionLogs, err := models.GetLogs(ctx, jardinId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, actionLogs)
}

func PostLog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var postedDTO dto.ActionLogDTO
	if err := c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	actionLog, err := postedDTO.ToActionLog()
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	id, err := actionLog.Save(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, map[string]string{"_id": id})
}

func PostLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var postedDTO []dto.ActionLogDTO
	if err := c.BindJSON(&postedDTO); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	logActions := make([]models.ActionLog, 0, len(postedDTO))
	for _, v := range postedDTO {
		a, err := v.ToActionLog()
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}
		logActions = append(logActions, a)
	}

	updatedLogs, err := models.SaveLogs(ctx, logActions)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("{updated: %d}", updatedLogs))
}

func DeleteLog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	a, err := models.GetLog(ctx, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err.Error())
		return
	}

	if err := models.DeleteLog(ctx, id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	DeleteBucketObjects(c, a.Photos)
	c.IndentedJSON(http.StatusOK, map[string]string{"_id": id})
}

func GetTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tags, err := models.GetTags(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, tags)
}

func GetLieux(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lieux, err := models.GetLieux(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, lieux)
}

func UpdateLogsSetGarden(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	valueStr := c.Query("value")
	if valueStr == "" {
		c.IndentedJSON(http.StatusBadRequest, "missing 'value' query parameter")
		return
	}

	nbModified, err := models.UpdateLogsSetGarden(ctx, valueStr)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, fmt.Sprintf("{updated: %d}", nbModified))
}
