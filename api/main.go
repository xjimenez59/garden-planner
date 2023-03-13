package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getLogs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, actionLogDemoData)
}

func main() {
	router := gin.Default()
	router.GET("/logs", getLogs)

	router.Run("localhost:8081")
}
