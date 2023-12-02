package main

import (
	"artika/api/user"
	"artika/client/template/pages"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionValidateArgs struct {
	SessionID string `json:"session-id"`
}

func routeSessionValidate(ctx *gin.Context) {
	var args SessionValidateArgs

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isSessionValid, err := user.IsSessionValid(args.SessionID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"is-session-valid": isSessionValid,
	})
}

type SessionCreateRequestBody struct {
	JWT string `json:"jwt"`
}

func routeSessionCreate(ctx *gin.Context) {
	var args SessionCreateRequestBody

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userInfo, err := user.DecodeJWT(args.JWT)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userSession, err := user.CreateSession(userInfo)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"session-id": userSession.UniqueID,
	})
}

func routeIndex(ctx *gin.Context) {
	component := pages.Index()
	component.Render(context.Background(), ctx.Writer)
}

type SessionArgs struct {
	SessionID string `json:"SessionID"`
}

func routeComponentAccount(ctx *gin.Context) {
	var args SessionArgs

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// _, err := user.IsSessionValid(args.SessionID)
	// if err != nil {
	// 	ctx.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }
}

func addAPIHeaders() gin.HandlerFunc {
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
	router.Use(addAPIHeaders())
	router.POST("/api/session/validate", routeSessionValidate)
	router.POST("/api/session/create", routeSessionCreate)

	router.Static("js", "./client/js")
	router.GET("/", routeIndex)
	router.GET("/component/account", routeComponentAccount)
	router.Run("localhost:3000")
}
