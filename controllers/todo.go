package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	data := new(models.Todo)

	err := c.BodyParser(&data)

	// if error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	data.ID = nil
	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	result, err := todoCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert todo",
			"error":   err,
		})
	}

	// get the inserted data
	todo := &models.Todo{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})
}

// get all todos
func GetTodos(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	// Query to filter
	query := bson.D{{}}

	cursor, err := todoCollection.Find(c.Context(), query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	var todos []models.Todo = make([]models.Todo, 0)

	// iterate the cursor and decode each item into a Todo
	err = cursor.All(c.Context(), &todos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todos": todos,
		},
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	// find parameter
	paramID := c.Params("id")

	// convert parameterID to objectId
	id, err := primitive.ObjectIDFromHex(paramID)

	// if parameter cannot parse
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	// var data Request
	data := new(models.Todo)
	err = c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	// updateData
	var dataToUpdate bson.D

	if data.Title != nil {
		// todo.Title = *data.Title
		dataToUpdate = append(dataToUpdate, bson.E{Key: "title", Value: data.Title})
	}

	if data.Completed != nil {
		// todo.Completed = *data.Completed
		dataToUpdate = append(dataToUpdate, bson.E{Key: "completed", Value: data.Completed})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	// update
	err = todoCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update todo",
			"error":   err,
		})
	}

	// get updated data
	todo := &models.Todo{}

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	// get param
	paramID := c.Params("id")

	// convert parameter to object id
	id, err := primitive.ObjectIDFromHex(paramID)

	// if parameter cannot parse
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	// find and delete todo
	query := bson.D{{Key: "_id", Value: id}}

	err = todoCollection.FindOneAndDelete(c.Context(), query).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete todo",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// sql filter
// elasticsearch
// kibana
// logstash
func GetTodo(id string) (*models.Todo, error) {
	// Kết nối đến Elasticsearch
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}

	// Truy vấn Elasticsearch để lấy todo dựa trên ID
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}

	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("todos"),
		esClient.Search.WithBody(esutil.NewJSONReader(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("Elasticsearch search request failed: %s", err)
	}
	defer res.Body.Close()

	// Xử lý kết quả tìm kiếm
	if res.IsError() {
		log.Fatalf("Elasticsearch search response error: %s", res.String())
	}

	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Fatalf("Failed to decode Elasticsearch response: %s", err)
	}

	hits := resBody["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		var id string
		if idVal, ok := source["id"].(string); ok {
			id = idVal
		}
		var title string
		if titleVal, ok := source["title"].(string); ok {
			title = titleVal
		}

		todo := models.Todo{
			ID:        &id,
			Title:     &title,
			Completed: new(bool),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		return &todo, nil
	}
	return nil, fmt.Errorf("Todo not found")
}
func GetTodoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	todo, err := GetTodo(id)
	if err != nil {
		return err
	}

	return c.JSON(todo)
}
