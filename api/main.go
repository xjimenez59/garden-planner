package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"garden-planner/api/config"
	"garden-planner/api/controllers"
)

func main() {

	config.ConnectDatabase()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/logs", controllers.GetLogs)
	router.POST("/logs", controllers.PostLogs)

	router.GET("/gardens", controllers.GetGardens)
	router.GET("/legumes", controllers.GetLegumes)

	router.GET("/tags", controllers.GetTags)

	router.Run("localhost:8081")
}
