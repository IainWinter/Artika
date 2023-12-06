package main

import (
	"artika/api/user"
	"artika/client/template/pages"
	"artika/client/template/prop"
	"artika/client/template/view"
	"context"
	"net/http"

	"github.com/a-h/templ"
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

// May want to split this up
// as it is not the api but the html generation
// For now I don't want 2 servers running, so just keep it here

// If the response should be a full page, or just its contents
func isRequestSPA(ctx *gin.Context) bool {
	hxRequestHeader := ctx.Request.Header.Get("Hx-Request")
	return hxRequestHeader == "true"
}

func getViewProps(ctx *gin.Context) (prop.ViewProps, error) {
	sessionIDCookie, err := ctx.Request.Cookie("SessionID")

	// If there is no cookie, that is not an error
	// Return an empty view props
	if err == http.ErrNoCookie {
		return prop.ViewProps{}, nil
	}

	return prop.GetViewPropsFromSessionID(sessionIDCookie.Value)
}

func returnDesktop(ctx *gin.Context, viewProps prop.ViewProps, component templ.Component) {
	if !isRequestSPA(ctx) {
		component = view.Desktop(viewProps, component)
	}

	// If the session is not valid, delete the cookie
	if !viewProps.IsSessionValid {
		ctx.Header("Set-Cookie", "SessionID=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT")
	}

	component.Render(context.Background(), ctx.Writer)
}

func routeIndex(ctx *gin.Context) {
	viewProps, err := getViewProps(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	returnDesktop(ctx, viewProps, pages.Index())
}

func routeAccount(ctx *gin.Context) {
	viewProps, err := getViewProps(ctx)
	if err != nil || !viewProps.IsSessionValid {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	returnDesktop(ctx, viewProps, pages.Account())
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
	router.GET("/account", routeAccount)

	router.Run("localhost:3000")
}
