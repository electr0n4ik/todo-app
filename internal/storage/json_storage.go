package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"todo-app/internal/todo"
)

// LoadJSON возвращает список задач из указанного файла и сохраняет их в tasks.json, перезаписав файл.
// Если указанного файла не обнаружено, то создает его и сообщает об этом. Работает с файлами json.
func LoadJSON(fileName string) ([]todo.Task, error) {
	filePath := filepath.Join(".", fileName)
	tasks := []todo.Task{}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		if err := SaveJSON(filePath, tasks); err != nil {
			return tasks, err
		}
	} else {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return tasks, err
		}
		if err := json.Unmarshal(content, &tasks); err != nil {
			return tasks, err
		}
		if err := SaveJSON("tasks.json", tasks); err != nil {
			return tasks, err
		}
	}

	return tasks, nil
}

// SaveJSON сохраняет список задач в указанном файле. Сообщает, если файл отсутствует.
// Проверяет наличие файла и наполняет его переданным списком задач. Сообщает, если файл не обнаружен.
// Файл имеет структуру json.
func SaveJSON(fileName string, tasks []todo.Task) error {
	filePath := filepath.Join(".", fileName)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Файл %s не обнаружен и будет создан!\n", fileName)
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			return err
		}
	}

	tasksJson, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, tasksJson, 0644); err != nil {
		return err
	}
	return nil
}
