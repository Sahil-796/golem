package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/Sahil-796/golem/server/pkg/balancer"
	// "fmt"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.Run("localhost:8080")
}
