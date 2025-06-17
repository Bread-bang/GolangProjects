package todo

import (
	"encoding/json"
	"os"
)

type Todo struct {
	Task string `json:"task"`
	Done bool `json:"done"`
}

func LoadTodos(fileName string) ([]Todo, error) {
	var todoList []Todo

	file, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return todoList, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&todoList)

	return todoList, nil
}

func SaveTodos(fileName string, todos []Todo) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(todos)
}