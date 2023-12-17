package main

import (
	"artika/api"
	"artika/client"

	"github.com/gin-gonic/gin"
)

// If the request is for the API, add headers
// This was mainly for the sign in with google request
func addHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.URL.Path) < 5 || c.Request.URL.Path[:4] != "/api" {
			c.Next()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	}
}

func main() {
	router := gin.Default()

	router.Use(addHeaders())

	api.AddRoutes(router)
	client.AddRoutes(router)

	router.Run("localhost:3000")
}
