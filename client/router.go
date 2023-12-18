package client

import (
	"artika/api/data"
	"artika/api/user"
	"artika/client/template/components"
	"artika/client/template/pages"
	"artika/client/template/prop"
	"artika/client/template/view"
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

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
		// This is is the wrong way to do this
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	returnDesktop(ctx, viewProps, pages.Account(viewProps.UserInfo))
}

func routeUserDesignerEnable(ctx *gin.Context) {
	sessionIDCookie, err := ctx.Request.Cookie("SessionID")

	var success = false

	if err == nil {
		err = user.EnableUserAsDesignerFromSessionID(sessionIDCookie.Value)
		if err == nil {
			success = true
		}
	}

	var component = components.UserEnableDesignerSnippet(success)
	component.Render(context.Background(), ctx.Writer)
}

func routeUserInfoEdit(ctx *gin.Context) {
	sessionIDCookie, err := ctx.Request.Cookie("SessionID")

	var args = data.UserInfoUpdate{
		FirstName: ctx.PostForm("firstName"),
		LastName:  ctx.PostForm("lastName"),
		//PictureURI: "",
		//Email:     ctx.PostForm("email"),
		//Address:   ctx.PostForm("address"),
		//City:      ctx.PostForm("city"),
		//State:     ctx.PostForm("state"),
		//Zip:       ctx.PostForm("zip"),
		//Country:   ctx.PostForm("country"),
	}

	var success = false

	if err == nil {
		err = user.UpdateUserInfoFromSessionID(sessionIDCookie.Value, args)
		if err == nil {
			success = true
		}
	}

	var component = components.UserEnableDesignerSnippet(success)
	component.Render(context.Background(), ctx.Writer)
}

func AddRoutes(router *gin.Engine) {
	router.Static("js", "./client/js")
	router.Static("css", "./client/css")

	// Pages
	router.GET("/", routeIndex)
	router.GET("/account", routeAccount)

	// Snippets which also have functionality
	router.POST("/user/info", routeUserInfoEdit)
	router.POST("/user/enableDesigner", routeUserDesignerEnable)
}
