package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/models"
	"go.elastic.co/apm"
)

func InsertTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)
	if err := c.BodyParser(&todo); err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	currentTime := time.Now()
	todo.CreatedAt = currentTime
	todo.UpdatedAt = currentTime

	stmt, err := config.MI.DB.PrepareContext(c.Context(), "INSERT INTO todos (userid, title, completed, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()

		fmt.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(c.Context(), todo.UserID, todo.Title, todo.Completed, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create todo"})
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		panic(err.Error())
	}
	todo.ID = &insertedID

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func GetAllTodos(c *fiber.Ctx) error {
	_, err := config.MI.DB.ExecContext(c.Context(), "USE todos")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		panic(err.Error())
	}

	rows, err := config.MI.DB.QueryContext(c.Context(), "SELECT * FROM todos")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()

		panic(err.Error())
	}

	defer rows.Close()

	todos := []models.Todo{}

	for rows.Next() {
		var todo models.Todo
		var createdAt, updatedAt []uint8
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &createdAt, &updatedAt, &todo.UserID); err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error0"})
		}
		createdAtTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid createdAt format"})
		}
		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid updatedAt format"})
		}
		todo.CreatedAt = createdAtTime
		todo.UpdatedAt = updatedAtTime
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func RemoveTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	var todoID int
	fmt.Sscanf(id, "%d", &todoID)

	_, err := config.MI.DB.Exec("USE todos")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		panic(err.Error())
	}

	stmt, err := config.MI.DB.PrepareContext(c.Context(), "DELETE FROM todos WHERE id = ?")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete todo"})
	}
	defer stmt.Close()

	_, err = stmt.Exec(todoID)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete todo"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Todo deleted successfully"})
}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	var todoID int
	fmt.Sscanf(id, "%d", &todoID)
	todo := new(models.Todo)
	if err := c.BodyParser(&todo); err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := config.MI.DB.Exec("USE todos")
	if err != nil {
		panic(err.Error())
	}
	stmt, err := config.MI.DB.Prepare("UPDATE todos SET user_id = ?, title = ?, completed = ? WHERE id = ?")
	if err != nil {
		fmt.Printf("Fail to prepare update statement: %s", err)
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fail to update todo"})
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.UserID, todo.Title, todo.Completed, todoID)
	if err != nil {
		fmt.Printf("Fail to update todo: %s", err)
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fail to update todo"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Todo deleted successfully"})
}

func GetTodo(c *fiber.Ctx) error {

	id := c.Params("id")

	var todoID int
	var createdAt []uint8
	var updatedAt []uint8
	fmt.Sscanf(id, "%d", &todoID)
	todo := new(models.Todo)

	_, err := config.MI.DB.ExecContext(c.Context(), "USE todos")
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		panic(err.Error())
	}

	row, err := config.MI.DB.QueryContext(c.Context(), "SELECT id, user_id, title, completed, createdAt, updatedAt  FROM todos WHERE id = ?", todoID)
	if err != nil {
		fmt.Printf("Fail to prepare update statement: %s", err)
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fail to get todo"})
	}
	if row.Next() {
		err := row.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed, &createdAt, &updatedAt)
		if err != nil {
			fmt.Println("Error getting todo:", err)
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get todo"})
		}
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	}

	createdAtTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse createdAt"})
	}
	updatedAtTime, err := time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse updatedAt"})
	}
	todo.CreatedAt = createdAtTime
	todo.UpdatedAt = updatedAtTime

	return c.JSON(todo)
}

func GetTodosByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	userTodos := []models.Todo{}

	rows, err := config.MI.DB.QueryContext(c.Context(), "SELECT id, title, completed, createdAt, updatedAt, user_id FROM todos WHERE user_id = ?", userID)
	if err != nil {
		apm.CaptureError(c.Context(), err).Send()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error1"})
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		var createdAt, updatedAt []uint8
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &createdAt, &updatedAt, &todo.UserID); err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error2"})
		}
		createdAtTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid createdAt format"})
		}
		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			apm.CaptureError(c.Context(), err).Send()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid updatedAt format"})
		}
		todo.CreatedAt = createdAtTime
		todo.UpdatedAt = updatedAtTime

		userTodos = append(userTodos, todo)
	}

	return c.JSON(userTodos)
}
