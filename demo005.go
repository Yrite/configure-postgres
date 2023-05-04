package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Define a struct to hold data from the database
type MyTable struct {
	Col1 string `json:"col1"`
	Col2 string `json:"col2"`
	Col3 string `json:"col3"`
}

func main() {
	// Connect to Postgres
	db, err := sql.Open("postgres", "user=your_username password=your_password dbname=your_dbname sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute the SQL query to select data from the table
	rows, err := db.Query("SELECT col1, col2, col3 FROM my_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create an empty slice to store the resulting data
	var data []MyTable

	// Iterate over the rows of the result set and store data in slice
	for rows.Next() {
		var col1, col2, col3 string

		if err := rows.Scan(&col1, &col2, &col3); err != nil {
			log.Fatal(err)
		}

		data = append(data, MyTable{
			Col1: col1,
			Col2: col2,
			Col3: col3,
		})
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Marshal the data into JSON
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON data to a file
	if err := os.WriteFile("output.json", jsonData, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data has been written to file output.json.")
}
H1uzkyfsVrnF