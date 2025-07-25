package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-app/internal/storage"
	"todo-app/internal/todo"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("todo: пропущена команда")
		fmt.Println("По команде \"./todo --help\" можно получить дополнительную информацию.")
		os.Exit(1)
	}
	args := os.Args[2:]
	command := os.Args[1]
	tasks, err := storage.LoadJSON("tasks.json")
	errCheck(err, "Ошибка загрузки задач из файла tasks.json!")

	switch command {
	case "list":
		// Как лучше оформлять длинные строки с передачей аргументов?
		// Можно переносить в следующую строку по одному аргументу?
		strFilter := *getNewFlagSet("list", "filter", args, "all", "--filter=all/done/pending")
		beautifulTasks, err := json.MarshalIndent(todo.List(tasks, strFilter), "", "    ")
		errCheck(err, "Ошибка вывода задач!")
		fmt.Println(string(beautifulTasks))

	case "add":
		strDesc := *getNewFlagSet("add", "desc", args, "", "--desc=<значение>")
		tasks = todo.Add(tasks, strDesc)
		errCheck(storage.SaveJSON("tasks.json", tasks), "Ошибка сохранения новой задачи!")

	case "complete":
		strId := *getNewFlagSet("complete", "id", args, "", "--id=<значение>")
		intId, err := strconv.Atoi(strId)
		errCheck(err, "Ошибка преобразования строки в число!")
		tasks, err = todo.Complete(tasks, intId)
		errCheck(err, "Ошибка изменения статуса задачи!")
		errCheck(
			storage.SaveJSON("tasks.json", tasks),
			"Ошибка сохранения задач!",
		)

	case "delete":
		errCheck(err, "Ошибка загрузки задач из файла tasks.json!")
		strId := *getNewFlagSet("delete", "id", args, "", "--id=<значение>")
		intId, err := strconv.Atoi(strId)
		errCheck(err, "Ошибка преобразования строки в число!")
		tasks, err = todo.Delete(tasks, intId)
		errCheck(err, "Ошибка изменения статуса задачи!")
		errCheck(
			storage.SaveJSON("tasks.json", tasks),
			"Ошибка сохранения задач!",
		)

	case "load":
		fileName := *getNewFlagSet("load", "file", args, "", "--file=<имя_файла.формат>")
		switch strings.Split(fileName, ".")[1] {
		case "json":
			_, err := storage.LoadJSON(fileName)
			errCheck(err, "Ошибка выгрузки задач из файла!")
		case "csv":
			_, err := storage.LoadCSV(fileName)
			errCheck(err, "Ошибка выгрузки задач из файла!")
		}
	case "export":
		commandCmd := flag.NewFlagSet("export", flag.ExitOnError)
		formatFile := commandCmd.String("format", "json", "--format=\"json\"/\"csv\"")
		fileName := commandCmd.String("out", "tasksExport.json", "--out=имяФайла.форматФайла")
		errCheck(commandCmd.Parse(args), "Ошибка получения значений флагов!")
		tasks, err := storage.LoadJSON("tasks.json")
		errCheck(err, "Ошибка выгрузки задач из файла!")
		switch *formatFile {
		case "json":
			if strings.Split(*fileName, ".")[1] != "json" {
				errCheck(errors.New(""), "Передан неверный файл или формат экспорта!")
			}
			errCheck(storage.SaveJSON(*fileName, tasks), "Ошибка экспортирования задач!")
		case "csv":
			if strings.Split(*fileName, ".")[1] != "csv" {
				errCheck(errors.New(""), "Передан неверный файл или формат экспорта!")
			}
			errCheck(storage.SaveCSV(*fileName, tasks), "Ошибка экспортирования задач!")
		default:
			errCheck(errors.New(""), "Передан неверный формат экспорта!")
			fmt.Println("По команде \"./todo --help\" можно получить дополнительную информацию.")
		}

	case "--help":
		fmt.Println("Использование: ./todo [command] --flag=<значение>")

	case "--version":
		fmt.Println("0.1.0")

	default:
		fmt.Println("todo: неизвестная команда:", command)
		fmt.Println("По команде \"./todo --help\" можно получить дополнительную информацию.")
		// В реальных программах делают принт ошибки в os.Stderr, почему так?
	}
}

// getNewFlagSet позволяет прочитать и получить значение флага
func getNewFlagSet(command string,
	flagName string,
	args []string,
	def_str string,
	info_str string) *string {
	commandCmd := flag.NewFlagSet(command, flag.ExitOnError)
	flagStr := commandCmd.String(flagName, def_str, info_str)
	if err := commandCmd.Parse(args); err != nil {
		fmt.Println("Ошибка чтения флага!")
		fmt.Println("По команде \"./todo --help\" можно получить дополнительную информацию.")
		os.Exit(1)
	}

	if *flagStr == "" {
		fmt.Println("Необходимо передать значение флага:")
		fmt.Println("По команде \"./todo --help\" можно получить дополнительную информацию.")
		os.Exit(1)
	}

	return flagStr
}

// errCheck позволяет обрабатывать ошибку компактнее
func errCheck(err error, errStr string) {
	if err != nil {
		fmt.Println(errStr)
		os.Exit(1)
	}
}
