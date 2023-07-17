package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "todos"
	password = "secret"
	hostname = "localhost:3306"
	dbname   = "todo_db"
)

type MySQLInstance struct {
	DB *sql.DB
}

var MI MySQLInstance

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func ConnectDB() {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL")

	MI = MySQLInstance{
		DB: db,
	}
}
