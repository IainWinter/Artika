package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// type Todo struct {
// 	Id        int    `json:"id"`
// 	Item      string `json:"item"`
// 	Completed bool   `json:"completed"`
// }

// type TodoBuildData struct {
// 	Item string `json:"item"`
// }

// var todos = []Todo{
// 	{Id: 1, Item: "test", Completed: true},
// 	{Id: 2, Item: "test", Completed: true},
// }

// func route_GetTodos(context *gin.Context) {
// 	context.IndentedJSON(http.StatusOK, todos)
// }

// func route_AddTodo(context *gin.Context) {
// 	var args TodoBuildData

// 	if err := context.BindJSON(&args); err != nil {
// 		context.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	var todo Todo
// 	todo.Id = len(todos)
// 	todo.Item = args.Item
// 	todo.Completed = false

// 	todos = append(todos, todo)

// 	context.IndentedJSON(http.StatusOK, todo)
// }

// func todo_GetTodoById(id int) (*Todo, error) {
// 	for i, t := range todos {
// 		if t.Id == id {
// 			return &todos[i], nil
// 		}
// 	}

// 	return nil, errors.New("No Todo found")
// }

// func route_SetTodoItem(context *gin.Context) {
// 	idString := context.Param("item")

// 	id, err := strconv.Atoi(idString)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is not an int"})
// 		return
// 	}

// 	todo, err := todo_GetTodoById(id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
// 		return
// 	}

// 	context.IndentedJSON(http.StatusOK, todo)
// }

type SessionCreateRequestBody struct {
	JWT string `json:"jwt"`
}

type MyClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
}

func route_session_create(ctx *gin.Context) {
	var args SessionCreateRequestBody

	if err := ctx.BindJSON(&args); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	publicKey, err := getGoogleCerts()

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// pass your custom claims to the parser function
	token, err := jwt.ParseWithClaims(args.JWT, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(publicKey), nil
	})

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// type-assert `Claims` into a variable of the appropriate type
	myClaims := token.Claims.(*MyClaims)

	ctx.IndentedJSON(http.StatusOK, gin.H{"jwt": myClaims})
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
