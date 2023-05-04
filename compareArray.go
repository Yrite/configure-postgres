package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Чтение первого файла с данными
	data1, err := ioutil.ReadFile("./data1.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла 1: %s\n", err.Error())
	}

	// Чтение второго файла с данными
	data2, err := ioutil.ReadFile("./data2.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла 2: %s\n", err.Error())
	}

	// Декодирование JSON массивов
	var array1, array2 []interface{}
	if err := json.Unmarshal(data1, &array1); err != nil {
		log.Fatalf("Ошибка при декодировании файла 1: %s\n", err.Error())
	}

	if err := json.Unmarshal(data2, &array2); err != nil {
		log.Fatalf("Ошибка при декодировании файла 2: %s\n", err.Error())
	}

	// Поиск различий между массивами
	diff := findArrayDifference(array1, array2)

	// Запись различий в файл
	diffData, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		log.Fatalf("Ошибка при кодировании результатов: %s\n", err.Error())
	}

	if err := ioutil.WriteFile("./diff.json", diffData, 0644); err != nil {
		log.Fatalf("Ошибка при записи в файл: %s\n", err.Error())
	}

	fmt.Println("Различия между файлами записаны в 'diff.json'")
}

// Функция для поиска различий между двумя массивами
func findArrayDifference(a, b []interface{}) []interface{} {
	diff := make([]interface{}, 0)

	for _, itemA := range a {
		found := false

		for _, itemB := range b {
			if isEqual(itemA, itemB) {
				found = true
				break
			}
		}

		if !found {
			diff = append(diff, itemA)
		}
	}

	return diff
}

// Функция для проверки наличия эквивалентности между двумя JSON объектами
func isEqual(a, b interface{}) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return string(aBytes) == string(bBytes)
}
