package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	PORT := os.Getenv("PORT")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + PORT)
	fmt.Printf("App is running on port %v\n", PORT)
}
