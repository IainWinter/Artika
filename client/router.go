package client

import (
	"artika/api/data"
	"artika/api/user"
	"artika/client/template/components"
	"artika/client/template/pages"
	"artika/client/template/prop"
	"artika/client/template/view"
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

var SessionNotValidErr = errors.New("Session not valid")

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
		var empty = prop.ViewProps{}
		empty.Url = ctx.Request.URL.String()

		return empty, nil
	}

	viewProps, err := prop.GetViewPropsFromSessionID(sessionIDCookie.Value)
	viewProps.Url = ctx.Request.URL.String()

	return viewProps, err
}

func getVerifiedSessionID(ctx *gin.Context) (string, error) {
	sessionIDCookie, err := ctx.Request.Cookie("SessionID")
	if err != nil {
		return "", err
	}

	isVerified, err := user.IsSessionValid(sessionIDCookie.Value)
	if !isVerified {
		return "", SessionNotValidErr
	}

	return sessionIDCookie.Value, err
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
	if err != nil || !viewProps.IsSessionValid { // unsure of all copies of this
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	returnDesktop(ctx, viewProps, pages.Account(viewProps.UserInfo))
}

func routeCreateRequest(ctx *gin.Context) {
	viewProps, err := getViewProps(ctx)
	if err != nil || !viewProps.IsSessionValid {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	returnDesktop(ctx, viewProps, pages.CreateRequest())
}

func routeMyRequests(ctx *gin.Context) {
	viewProps, err := getViewProps(ctx)
	if err != nil || !viewProps.IsSessionValid {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	workRequests, err := user.GetAllWorkItemsForValidSessionID(viewProps.SessionID)

	returnDesktop(ctx, viewProps, pages.MyRequests(workRequests))
}

func routeUserDesignerEnable(ctx *gin.Context) {
	sessionIDCookie, err := ctx.Request.Cookie("SessionID")

	// verify session

	var success = false

	if err == nil {
		err = user.EnableUserAsDesignerForValidSessionID(sessionIDCookie.Value)
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
		err = user.UpdateUserInfoForValidSessionID(sessionIDCookie.Value, args)
		if err == nil {
			success = true
		}
	}

	var component = components.UserEnableDesignerSnippet(success)
	component.Render(context.Background(), ctx.Writer)
}

func routeWorkItemCreate(ctx *gin.Context) {
	sessionID, err := getVerifiedSessionID(ctx)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	file, header, err := ctx.Request.FormFile("test-file")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	picture, err := user.StorePictureForValidSessionID(sessionID, file, header)

	var success = false

	if err == nil {
		var args = data.WorkItemCreateInfo{
			Title:       ctx.PostForm("title"),
			Description: ctx.PostForm("description"),
			PictureIDs:  []string{picture.PictureID},
		}

		_, err := user.CreateWorkItemForValidSessionID(sessionID, args)
		if err == nil {
			success = true
		}
	}

	var component = components.UserEnableDesignerSnippet(success) // replace
	component.Render(context.Background(), ctx.Writer)
}

func routeWorkItemDelete(ctx *gin.Context) {
	sessionID, err := getVerifiedSessionID(ctx)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	workItemID := ctx.Param("workItemID")

	err = user.DeleteWorkItemForValidSessionID(sessionID, workItemID)

	var success = false

	if err == nil {
		success = true
	}

	var component = components.UserEnableDesignerSnippet(success) // replace
	component.Render(context.Background(), ctx.Writer)
}

func AddRoutes(router *gin.Engine) {
	router.Static("js", "./client/js")
	router.Static("css", "./client/css")

	// Pages
	router.GET("/", routeIndex)
	router.GET("/account", routeAccount)
	router.GET("/createRequest", routeCreateRequest)
	router.GET("/myRequests", routeMyRequests)

	// Snippets which also have functionality
	router.POST("/user/info", routeUserInfoEdit)
	router.POST("/user/enableDesigner", routeUserDesignerEnable)

	router.POST("/workItem", routeWorkItemCreate)
	router.DELETE("/workItem/:workItemID", routeWorkItemDelete)
}
