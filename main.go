package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"

	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/controller"
	"github.com/nguyendinhduy02112002/golang_todo_app/db"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmfiber"
)

func setupTodosAPI(app *fiber.App) {
	app.Put("api/todos/:id", controller.UpdateTodo)
	app.Get("api/todos/:id", controller.GetTodo)
	app.Get("api/todos", controller.GetAllTodos)
	app.Get("api/todos/user/:id", controller.GetTodosByUserID)
	app.Post("api/todos", controller.InsertTodo)
	app.Delete("api/todos/:id", controller.RemoveTodo)

}

func main() {
	tracer, err := apm.NewTracer("todos_service", "2.0.0")
	if err != nil {
		fmt.Println(err)
	}
	defer tracer.Close()

	config.NewElasticsearchClient()
	config.ConnectDB()
	app := fiber.New()

	db.MigrateDB()
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(tracer)))

	setupTodosAPI(app)

	defer tracer.Flush(nil)

	app.Listen(":3001")

}
