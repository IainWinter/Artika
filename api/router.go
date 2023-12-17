package api

import (
	"artika/api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeSessionCreate(ctx *gin.Context) {
	var args struct {
		JWT string `json:"JWT"`
	}

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

func routeSessionDelete(ctx *gin.Context) {
	var args struct {
		SessionID string `json:"SessionID"`
	}

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

func routeDesignersGetAllPublic(ctx *gin.Context) {
	designers, err := user.GetAllPublicDesigners()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"Designers": designers,
	})
}

func routeUserEnableAsDesigner(ctx *gin.Context) {
	var args struct {
		SessionID string `json:"SessionID"`
	}

	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = user.EnableUserAsDesignerFromSessionID(args.SessionID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{})
}

func AddRoutes(router *gin.Engine) {
	router.POST("/api/session", routeSessionCreate)
	router.DELETE("/api/session", routeSessionDelete)

	router.POST("/api/user/enableDesigner", routeUserEnableAsDesigner)

	router.GET("/api/designers", routeDesignersGetAllPublic)
}
