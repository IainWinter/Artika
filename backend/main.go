package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionCreateRequestBody struct {
	JWT string `json:"jwt"`
}

func route_session_create(ctx *gin.Context) {
	var args SessionCreateRequestBody

	if err := ctx.BindJSON(&args); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user_info_all, err := decode_jwt(args.JWT)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var name = user_info_all["name"].(string)
	var email = user_info_all["email"].(string)
	var picture = user_info_all["picture"].(string)

	var user_info = map[string]string{
		"name":    name,
		"email":   email,
		"picture": picture,
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"jwt": user_info})
}

func addAPIHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(addAPIHeaders())
	router.POST("/api/session/create", route_session_create)
	router.Run("localhost:3001")
}
