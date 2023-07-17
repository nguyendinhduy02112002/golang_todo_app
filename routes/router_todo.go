package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/controller"
)

func setRouter(route fiber.Router) {
	route.Get("", controller.GetTodos)
	route.Post("", controller.UpdateTodo)
}
