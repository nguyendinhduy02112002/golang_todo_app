package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/models"
)

func InsertTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)
	*todo.ID = "1"
	*todo.Title = "coding"
	*todo.Completed = false

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	stmt, err := config.MI.DB.Prepare("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Failed to prepare insert statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.ID, todo.Title, todo.Completed)
	if err != nil {
		fmt.Println("Failed to insert todo:", err)
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Todo created successfully"})
}
