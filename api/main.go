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
	//	router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "locale", ""},
		ExposeHeaders: []string{"Content-Length"},

		//	AllowCredentials: true,
		//	AllowOriginFunc: func(origin string) bool {		return true	},
	}))

	router.GET("/logs", controllers.GetLogs)
	router.POST("/logs", controllers.PostLogs)
	router.POST("/log", controllers.PostLog)
	router.DELETE("/log/:id", controllers.DeleteLog)
	router.PUT("/logs/garden", controllers.UpdateLogsSetGarden)

	router.GET("/gardens", controllers.GetGardens)
	router.POST("/garden", controllers.PostGarden)
	router.GET("/garden/:gardenId/logs", controllers.GetLogs)

	router.GET("/legumes", controllers.GetLegumes)

	router.GET("/tags", controllers.GetTags)
	router.GET("/lieux", controllers.GetLieux)

	router.POST("/photo", controllers.HandleFileUploadToBucket)
	router.DELETE("/photo/:id", controllers.DeleteBucketObject)

	router.GET("/recoltes", controllers.GetRecoltes)
	router.GET("/recoltes/lieux", controllers.GetRecoltesLieux)

	router.Run("0.0.0.0:8081")
}
