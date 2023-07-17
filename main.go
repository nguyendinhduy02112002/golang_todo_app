package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/controller"
	"github.com/nguyendinhduy02112002/golang_todo_app/db"
)

func main() {
	config.NewElasticsearchClient()
	config.ConnectDB()
	app := fiber.New()

	db.MigrateDB()

	app.Get("", controller.GetTodos)
	app.Put("/:id", controller.UpdateTodo)
	app.Get("/:id", controller.GetTodo)
	app.Post("/:id", controller.CreateTodo)
	app.Delete("/:id", controller.DeleteTodo)
	// // app.Post("/todos", controller.InsertTodo)

	app.Listen(":8080")

}
