package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DBIndex struct {
	Table        string `json:"table"`
	IndexName    string `json:"index_name"`
	ColumnNames  string `json:"column_names"`
	IsUnique     bool   `json:"is_unique"`
	IsPrimary    bool   `json:"is_primary"`
	IsPartial    bool   `json:"is_partial"`
	IsConcurrent bool   `json:"is_concurrent"`
}

func main() {
	// установка соединения с базой данных Postgres
	db, err := sql.Open("postgres", "user=имя_пользователя password=пароль dbname=имя_базы данных sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// выполнение SQL-запроса для получения индексов всех таблиц в базе данных
	rows, err := db.Query(`
        SELECT
            t.relname as table_name,
            ix.relname as index_name,
            array_to_string(array_agg(a.attname), ',') as column_names,
            ix.indisunique as is_unique,
            ix.indisprimary as is_primary,
            ix.indispartial as is_partial,
            ix.indisvalid as is_concurrent
        FROM
            pg_index as i
            JOIN pg_class as t ON i.indrelid = t.oid
            JOIN pg_class as ix ON i.indexrelid = ix.oid
            JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(i.indkey)
            JOIN pg_namespace nsp on nsp.oid = t.relnamespace
        WHERE nsp.nspname = 'tosca'
        GROUP BY table_name, index_name, is_unique, is_primary, is_partial, is_concurrent
    `)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// перебор результатов и создание словаря DBIndex-ов
	indices := make([]DBIndex, 0)
	for rows.Next() {
		index := DBIndex{}
		if err := rows.Scan(
			&index.Table,
			&index.IndexName,
			&index.ColumnNames,
			&index.IsUnique,
			&index.IsPrimary,
			&index.IsPartial,
			&index.IsConcurrent); err != nil {
			panic(err)
		}
		indices = append(indices, index)
	}

	// кодирование результата в формат JSON
	jsonOutput, err := json.Marshal(indices)
	if err != nil {
		panic(err)
	}

	// сохранение результата в файл
	outputFile, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	outputFile.Write(jsonOutput)

	// вывод полученных данных
	fmt.Println(string(jsonOutput))
}
