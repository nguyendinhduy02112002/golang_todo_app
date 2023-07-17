package db

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/nguyendinhduy02112002/golang_todo_app/config"
)

func MigrateDB() {
	_, err := config.MI.DB.Exec("USE todo_db")
	if err != nil {
		panic(err.Error())
	}

	rows, err := config.MI.DB.Query("SELECT id, title FROM todos")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		err := rows.Scan(&id, &title)
		if err != nil {
			log.Println(err)
			continue
		}

		// Gửi dữ liệu đến Elasticsearch
		doc := map[string]interface{}{
			"id":    id,
			"title": title,
		}

		_, err = config.EI.Client.Index("todos", esutil.NewJSONReader(&doc), config.EI.Client.Index.WithDocumentID(fmt.Sprintf("%d", id)))
		if err != nil {
			log.Println(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
