package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"garden-planner/meteo/config"
	"garden-planner/meteo/controllers"
)

func main() {

	config.ConnectDatabase()

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

	router.GET("/meteo/:site/:date", controllers.GetMeteo)

	router.Run("0.0.0.0:8082")
}
