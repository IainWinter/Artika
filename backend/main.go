package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

type TodoBuildData struct {
	Item string `json:"item"`
}

var todos = []Todo{
	{Id: 1, Item: "test", Completed: true},
	{Id: 2, Item: "test", Completed: true},
}

func route_GetTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func route_AddTodo(context *gin.Context) {
	var args TodoBuildData

	if err := context.BindJSON(&args); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var todo Todo
	todo.Id = len(todos)
	todo.Item = args.Item
	todo.Completed = false

	todos = append(todos, todo)

	context.IndentedJSON(http.StatusOK, todo)
}

func todo_GetTodoById(id int) (*Todo, error) {
	for i, t := range todos {
		if t.Id == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("No Todo found")
}

func route_SetTodoItem(context *gin.Context) {
	idString := context.Param("item")

	id, err := strconv.Atoi(idString)

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is not an int"})
		return
	}

	todo, err := todo_GetTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	todo.Item =

		context.IndentedJSON(http.StatusOK, todo)
}

func main() {
	fmt.Println("Hello")
	router := gin.Default()
	router.GET("/todos", route_GetTodos)
	router.POST("/todo", route_AddTodo)
	router.PATCH("/todo/item/:item", route_SetTodoItem)

	router.Run("localhost:9090")
}
