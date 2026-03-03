package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"garden-planner/meteo/config"
	"garden-planner/meteo/controllers"
)

func main() {

	config.ConnectDatabase()
	defer config.CloseDatabase()

	router := gin.Default()
	//	router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "locale", ""},
		ExposeHeaders: []string{"Content-Length"},

		//	AllowCredentials: true,
		//	AllowOriginFunc: func(origin string) bool {		return true	},
	}))

	router.GET("/meteo/infoclimat/:site/:date", controllers.GetMeteo)
	router.GET("/meteo/accuweather/:location/past24h", controllers.Accuweather_GetCurrent24H)
	router.GET("/meteo/accuweather/:location/past24h/import", controllers.Accuweather_ImportCurrent24H)
	router.GET("/meteo/accuweather/location/search", controllers.Accuweather_SearchLocations)
	router.GET("/meteo/accuweather/location/import", controllers.Accuweather_ImportLocations)

	router.GET("/meteofrance/quotidien", controllers.MF_CommandeQuotidienne)
	router.GET("/meteofrance/resultats", controllers.MF_GetResultats)
	router.GET("/meteo", controllers.MF_GetMeteo)

	router.Run("0.0.0.0:8082")

}
