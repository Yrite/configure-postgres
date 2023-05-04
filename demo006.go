package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Result struct {
	StringValue string `json:"stringValue"`
	IntValue    int    `json:"intValue"`
}

func main() {
	// Подключение к базе данных Postgres
	db, err := sql.Open("postgres", "user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения с базой данных: %s", err.Error())
	}
	defer db.Close()

	// Выполнение SQL-запросов и форматирование результатов в JSON
	results := []Result{}
	query1 := "SELECT column1, column2 FROM mytable WHERE column3 = $1"
	rows, err := db.Query(query1, "condition")
	if err != nil {
		log.Fatalf("Ошибка при выполнении SQL-запроса: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var col1 string
		var col2 int
		err := rows.Scan(&col1, &col2)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строк результата: %s", err.Error())
		}
		results = append(results, Result{StringValue: col1, IntValue: col2})
	}

	query2 := "SELECT COUNT(*) FROM mytable WHERE column3 = $1"
	var count int
	err = db.QueryRow(query2, "condition").Scan(&count)
	if err != nil {
		log.Fatalf("Ошибка при выполнении SQL-запроса: %s", err.Error())
	}
	results = append(results, Result{StringValue: "count", IntValue: count})

	jsonString, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Ошибка при форматировании результата в JSON: %s", err.Error())
	}
	fmt.Println(string(jsonString))
}
