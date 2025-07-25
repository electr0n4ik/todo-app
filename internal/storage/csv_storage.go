package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"todo-app/internal/todo"
)

// LoadCSV возвращает список задач из указанного файла и сохраняет их в tasks.json, перезаписав файл.
// Если указанного файла не обнаружено, то создает его и сообщает об этом. Работает с файлами csv.
func LoadCSV(fileName string) ([]todo.Task, error) {
	filePath := filepath.Join(".", fileName)
	tasks := []todo.Task{}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		if err := SaveCSV(filePath, tasks); err != nil {
			return tasks, err
		}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return tasks, err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return tasks, err
	}

	for i, v := range records {
		if i == 0 {
			continue
		}
		if len(v) != 3 {
			fmt.Printf("Во время чтения строк в файле %s обнаружена ошибка.", fileName)
			fmt.Printf("Ошибка в строке %d", i)
			fmt.Println("Не хватает данных для выгрузки задачи.")
			continue
		}
		idTask, err := strconv.Atoi(v[0])
		if err != nil {
			return tasks, err
		}
		doneTask, err := strconv.ParseBool(v[2])
		if err != nil {
			return tasks, err
		}
		tempTask := todo.Task{
			ID:          idTask,
			Description: v[1],
			Done:        doneTask,
		}
		tasks = append(tasks, tempTask)
	}

	if err := SaveJSON("tasks.json", tasks); err != nil {
		return tasks, err
	}

	return tasks, nil
}

// SaveCSV сохраняет список задач в указанном файле. Сообщает, если файл отсутствует.
// Всегда создает/пересоздает файл и наполняет его переданным списком задач.
// Файл имеет заголовки ID, Description, Done.
func SaveCSV(fileName string, tasks []todo.Task) error {
	filePath := filepath.Join(".", fileName)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Файл %s не обнаружен и будет создан!\n", fileName)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write([]string{"ID", "Description", "Done"}); err != nil {
		return err
	}

	for _, v := range tasks {
		strSliceTask := []string{
			strconv.Itoa(v.ID),
			v.Description,
			strconv.FormatBool(v.Done),
		}
		if err := csvWriter.Write(strSliceTask); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return err
	}

	return nil
}
