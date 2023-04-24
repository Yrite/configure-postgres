package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	// Подключение к базе данных
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем от пользователя данные для создания базы данных, схемы и пользователя
	var dbName string
	fmt.Print("Введите название новой базы данных: ")
	fmt.Scanln(&dbName)

	var schemaName string
	fmt.Print("Введите название новой схемы: ")
	fmt.Scanln(&schemaName)

	var userName string
	fmt.Print("Введите имя нового пользователя: ")
	fmt.Scanln(&userName)

	var userPassword string
	fmt.Print("Введите пароль для нового пользователя: ")
	fmt.Scanln(&userPassword)

	// Создание новой базы данных
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Создана новая база данных: %s", dbName)

	// Подключение к созданной базе данных
	db.Close()

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание новой схемы
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Создана новая схема: %s", schemaName)

	// Создание нового пользователя
	_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", userName, userPassword))
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

	// Устанавливаем путь поиска по умолчанию в созданной схеме для созданного пользователя
	_, err = db.Exec(fmt.Sprintf("ALTER USER %s SET search_path = %s", userName, schemaName))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Пользователю %s установлен путь поиска по умолчанию в схему %s", userName, schemaName)
}
