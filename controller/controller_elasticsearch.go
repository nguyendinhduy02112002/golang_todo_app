package controller

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gofiber/fiber/v2"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
	"github.com/nguyendinhduy02112002/golang_todo_app/models"
)

func GetTodos(c *fiber.Ctx) error {
	searchRequest := esapi.SearchRequest{
		Index: []string{"todos"},
		Body: strings.NewReader(`{
			"query": {
				"match_all": {}
			}
		}`),
	}

	res, err := searchRequest.Do(context.Background(), config.EI.Client)

	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer res.Body.Close()
	c.SendString(res.String())
	return err
}

func UpdateTodo(c *fiber.Ctx) error {

	documentID := c.Params("id")

	data := new(models.Todo)
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("cannot parse JSON: %s", err)
	}
	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error converting data to JSON: %s", err)
	}
	putRequest := esapi.IndexRequest{
		Index:      "todos",
		DocumentID: documentID,
		Body:       strings.NewReader(string(dataJSON)),
		Refresh:    "true",
	}

	res, err := putRequest.Do(context.Background(), config.EI.Client)
	if err != nil {
		log.Fatalf("failed to put document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("failed to put document: %s", res.String())
	}
	c.SendString(res.String())

	return c.SendStatus(fiber.StatusOK)
}

func GetTodo(c *fiber.Ctx) error {
	documentID := c.Params("id")
	getRequest := esapi.GetRequest{
		Index:      "todos",
		DocumentID: documentID,
	}

	res, err := getRequest.Do(context.Background(), config.EI.Client)
	if err != nil {
		log.Fatalf("err get data: %s", err)
	}

	defer res.Body.Close()
	c.SendString(res.String())
	return c.SendStatus(fiber.StatusOK)
}

func CreateTodo(c *fiber.Ctx) error {
	documentID := c.Params("id")
	data := new(models.Todo)

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("cannot parse JSON: %s", err)
	}

	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error converting data to JSON: %s", err)
	}

	createRequest := esapi.CreateRequest{
		Index:      "todos",
		DocumentID: documentID,
		Body:       strings.NewReader(string(dataJSON)),
	}

	res, err := createRequest.Do(context.Background(), config.EI.Client)
	if err != nil {
		log.Fatalf("err creating data: %s", err)
	}
	defer res.Body.Close()
	return c.SendStatus(fiber.StatusOK)
}

func DeleteTodo(c *fiber.Ctx) error {
	documentID := c.Params("id")
	deleteTodo := esapi.DeleteRequest{
		Index:      "todos",
		DocumentID: documentID,
	}

	res, err := deleteTodo.Do(context.Background(), config.EI.Client)
	if err != nil {
		log.Fatalf("err deleting data: %s", err)
	}
	defer res.Body.Close()
	return c.SendStatus(fiber.StatusOK)
}
