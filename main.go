package main

import (
	"github.com/Sahil-796/golem/config"
	"github.com/gin-gonic/gin"
	// "github.com/Sahil-796/golem/server/pkg/balancer"
	// "fmt"
)

func main() {
	
	config.LoadConfig()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello dumb fuck",
		})
	})

	router.Run("localhost:8080")
}
