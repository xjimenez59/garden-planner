package controllers

import (
	"context"
	"garden-planner/meteo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Accuweather_GetCurrent24H(c *gin.Context) {
	location := c.Param("location")

	data, err := models.Accuweather_GetCurrent24H(location)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func Accuweather_ImportCurrent24H(c *gin.Context) {

	//-- retrieve the mongo id of the queried location
	locationKey := c.Param("location")
	ctx := context.Background()
	location, err := models.Get_Location_byKey(ctx, locationKey)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}

	//-- get  data from Accuweather
	data, err := models.Accuweather_GetCurrent24H(location.Key)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// transform and save data
	for _, v := range data {
		var wh models.WeatherHourlyData
		wh.FromAccuweatherCurrent(v)
		wh.LocationID = location.ID
		if err := wh.Save(ctx); err != nil {
			log.Fatal(err)
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]string{"msg": "Successfully imported"})
}

func Accuweather_SearchLocations(c *gin.Context) {
	cp := c.Query("cp")

	data, err := models.Accuweather_GetLocations_from_CP(cp)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func Accuweather_ImportLocations(c *gin.Context) {
	cp := c.Query("cp")
	ctx := context.Background()

	//-- get data from Accuweather API
	data, err := models.Accuweather_GetLocations_from_CP(cp)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	//-- transform data into  our database format
	locations := make([]models.WeatherLocation, len(data))
	for i, v := range data {
		locations[i], err = models.FromAccuweatherLocation(v)
	}

	//-- save data into database
	for _, wl := range locations {
		if err := wl.Save(ctx); err != nil {
			log.Fatal(err)
		}
	}

	c.IndentedJSON(http.StatusOK, locations)
}
