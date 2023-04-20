package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgres://<user>:<password>@<host>:<port>/<database_name>?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание новой базы данных
	dbName := "newdatabase"
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Создана новая база данных: %s", dbName)

	// Создание новой схемы
	schemaName := "newschema"
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Создана новая схема: %s", schemaName)

	// Создание нового пользователя
	userName := "newuser"
	password := "newpassword"
	_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", userName, password))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Создан новый пользователь: %s", userName)

	// Дать созданному пользователю права на созданную схему
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON SCHEMA %s TO %s", schemaName, userName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Пользователю %s даны права на схему %s", userName, schemaName)
}
