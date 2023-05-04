package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Чтение первого файла
	file1, err := os.Open("file1.json")
	if err != nil {
		log.Fatalf("Ошибка при открытии файла 1: %s", err.Error())
	}
	defer file1.Close()

	bytes1, err := ioutil.ReadAll(file1)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла 1: %s", err.Error())
	}

	// Чтение второго файла
	file2, err := os.Open("file2.json")
	if err != nil {
		log.Fatalf("Ошибка при открытии файла 2: %s", err.Error())
	}
	defer file2.Close()

	bytes2, err := ioutil.ReadAll(file2)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла 2: %s", err.Error())
	}

	// Парсинг объектов JSON из файлов
	var obj1, obj2 map[string]interface{}

	if err := json.Unmarshal(bytes1, &obj1); err != nil {
		log.Fatalf("Ошибка при парсинге файла 1: %s", err.Error())
	}

	if err := json.Unmarshal(bytes2, &obj2); err != nil {
		log.Fatalf("Ошибка при парсинге файла 2: %s", err.Error())
	}

	// Сравнение объектов JSON и запись различий в файл
	diffObj := make(map[string]interface{})
	for k, v := range obj1 {
		if val, ok := obj2[k]; !ok || !isEqual(v, val) {
			diffObj[k] = v
		}
	}

	if err := writeJSONToFile(diffObj, "diff.json"); err != nil {
		log.Fatalf("Ошибка при записи файла: %s", err.Error())
	}
}

// Функция проверки равенства двух объектов JSON
func isEqual(a, b interface{}) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return string(aBytes) == string(bBytes)
}

// Функция записи объекта JSON в файл
func writeJSONToFile(obj map[string]interface{}, filename string) error {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, jsonBytes, 0644); err != nil {
		return err
	}

	fmt.Printf("Результат записан в файл '%s'\n", filename)
	return nil
}
