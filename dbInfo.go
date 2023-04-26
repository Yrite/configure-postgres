package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

// define JSON structure
type Table struct {
	SchemaName string `json:"schema"`
	TableName  string `json:"table"`
	UserName   string `json:"user"`
	TableSpace string `json:"tablespace"`
	DbList	   string `json:"dblist"`
}

func main() {
	port := 5432
	user := "db_user"
	dbname := "your_db"
	password := "your_pass"
	host := "localhost"

	// open db connection
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// get db tables 
	rows, err := db.Query(`SELECT table_name FROM information_schema.tables 
	WHERE table_schema NOT IN ('information_schema', 'pg_catalog') 
	AND table_type = 'BASE TABLE';`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}

	// get db schemas 
	rows, err = db.Query("SELECT schema_name FROM information_schema.schemata;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
 
	schemas := make([]string, 0)
	for rows.Next() {
		var schema string
		err := rows.Scan(&schema)
		if err != nil {
			log.Fatal(err)
		}
		schemas = append(schemas, schema)
	}

	// get db users
	rows, err = db.Query("SELECT usename FROM pg_user;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make([]string, 0)
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	// get pg space list 
	rows, err = db.Query("SELECT spcname FROM pg_tablespace;") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tableSpaces := make([]string, 0)
	for rows.Next() {
		var tableSpace string
		err := rows.Scan(&tableSpace)
		if err != nil {
			log.Fatal(err)
		}
		tableSpaces = append(tableSpaces, tableSpace)
	}

	rows, err = db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dbList := make([]string, 0)
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			log.Fatal(err)
		}
		dbList = append(dbList, db)
	}

	// create JSON struct
	data := map[string]interface{}{
		"schema"		: schemas,
		"db_tables"		: tables,
		"user"			: users,
		"tablespace"	: tableSpaces,
		"database"		: dbList, 
	}

	// convert struct to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	
	// create JSON db report 
	file, err := os.Create("db_info.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}
