package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/controllers" // replace
)

func TodoRoute(route fiber.Router) {
	route.Get("", controllers.GetTodos)
	route.Post("", controllers.CreateTodo)
	route.Delete("/:id", controllers.DeleteTodo)
	route.Put("/:id", controllers.UpdateTodo)
	route.Get("/:id", controllers.GetTodoHandler)

}

// ELK stack search engine
