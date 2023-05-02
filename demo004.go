package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Column struct {
	Name              string         `json:"column_name"`
	DataType          string         `json:"data_type"`
	IsNullable        string         `json:"is_nullable"`
	MaxLength         sql.NullInt64  `json:"character_maximum_length"`
	NumericPrecision  sql.NullInt64  `json:"numeric_precision"`
	DateTimePrecision sql.NullInt64  `json:"datetime_precision"`
	DefaultVal        sql.NullString `json:"column_default"`
}

type Table struct {
	Name    string   `json:"table_name"`
	Columns []Column `json:"columns"`
}

func main() {
	connStr := "user=your_username password=your_password dbname=your_dbname sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
        SELECT 
            table_name, 
            column_name, 
            data_type, 
            is_nullable,
            character_maximum_length,
            numeric_precision,
            datetime_precision,
            column_default
        FROM 
            information_schema.columns
        WHERE 
            table_schema = 'public' -- replace 'public' with name of schema you are interested in
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []Table
	var prevTable string // keep track of previous table name
	var currentTable *Table

	for rows.Next() {
		var column Column
		var tableName string

		err := rows.Scan(&tableName, &column.Name, &column.DataType, &column.IsNullable, &column.MaxLength, &column.NumericPrecision, &column.DateTimePrecision, &column.DefaultVal)
		if err != nil {
			log.Fatal(err)
		}

		// if table changed, create new Table instance
		if prevTable != tableName {
			currentTable = &Table{Name: tableName}
			tables = append(tables, *currentTable)
			prevTable = tableName
		}

		currentTable.Columns = append(currentTable.Columns, column)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// marshal tables to JSON
	jsonTables, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// write JSON to file
	file, err := os.Create("tables.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonTables)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables JSON written to file")
}
