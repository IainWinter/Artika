package main

import (
	"artika/api/user"
	"artika/client/template/pages"
	"artika/client/template/prop"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionCreateRequestBody struct {
	JWT string `json:"JWT"`
}

func routeSessionCreate(ctx *gin.Context) {
	var args SessionCreateRequestBody

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userSession, err := user.RegisterJWT(args.JWT)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"SessionID": userSession.UniqueID,
	})
}

type SessionDeleteRequestBody struct {
	SessionID string `json:"SessionID"`
}

func routeSessionDelete(ctx *gin.Context) {
	var args SessionDeleteRequestBody

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = user.DeleteSession(args.SessionID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{})
}

func routeIndex(ctx *gin.Context) {
	var pageProps prop.IndexProps

	sessionIDCookie, err := ctx.Request.Cookie("SessionID")
	if err != http.ErrNoCookie {
		pageProps, err = prop.GetIndexPagePropsFromSessionID(sessionIDCookie.Value)
	}

	var component = pages.Index(pageProps)
	component.Render(context.Background(), ctx.Writer)
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

	router.POST("/api/session", routeSessionCreate)
	router.DELETE("/api/session", routeSessionDelete)

	router.Static("js", "./client/js")
	router.Static("css", "./client/css")
	router.GET("/", routeIndex)
	router.Run("localhost:3000")
}
