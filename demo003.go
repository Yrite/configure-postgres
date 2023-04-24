package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	conStr := "postgres://user:password@localhost/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Назначение пользователя владельцем схемы
	_, err = db.Exec("ALTER SCHEMA my_schema OWNER TO my_user")
	if err != nil {
		panic(err)
	}

	// Назначение пользователя владельцем базы данных
	_, err = db.Exec("ALTER DATABASE mydatabase OWNER TO my_user")
	if err != nil {
		panic(err)
	}

	// Дать пользователю все привилегии в базе данных
	_, err = db.Exec("GRANT ALL PRIVILEGES ON DATABASE mydatabase TO my_user")
	if err != nil {
		panic(err)
	}

	// Перемещение всех таблиц базы данных из схемы postgres в пользовательскую схему
	_, err = db.Exec("ALTER TABLE my_schema.* SET SCHEMA my_user")
	if err != nil {
		panic(err)
	}
}
