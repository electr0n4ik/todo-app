package todo

import (
	"errors"
)

// Add добавляет задачу с переданным описанием в переданный список.
// Возвращает обновленный список.
func Add(tasks []Task, desc string) []Task {
	var id int = 0
	if len(tasks) > 0 {
		for _, v := range tasks {
			if v.ID > id {
				id = v.ID
			}
		}
	} else {
		id = 0
	}
	newTask := Task{id + 1, desc, false}

	return append(tasks, newTask)
}

// List возвращает список задач с указанным фильтром.
// Если передан фильтр all, то возвращает все задачи.
// Если передан фильтр done, то возвращает выполненные задачи.
// Если передан фильтр pending, то возвращает задачи в процессе выполнения.
func List(tasks []Task, filter string) []Task {
	tasksFilter := make([]Task, 0)
	switch filter {
	case "all":
		tasksFilter = tasks
	case "done":
		for _, v := range tasks {
			if v.Done {
				tasksFilter = append(tasksFilter, v)
			}
		}
	case "pending":
		for _, v := range tasks {
			if !v.Done {
				tasksFilter = append(tasksFilter, v)
			}
		}
	}
	return tasksFilter
}

// Complete меняет статус задачи на выполненную по переданному id.
// После изменения статуса возвращает обновленный список задач.
func Complete(tasks []Task, id int) ([]Task, error) {
	if len(tasks) > 0 {
		for i, v := range tasks {
			if v.ID == id {
				tasks[i].Done = true
				return tasks, nil
			}
		}

	}
	return tasks, errors.New("Нет задач на выполнение!")
}

// Delete удаляет задачу по переданному id.
// После удаления возвращает обновленный список задач.
func Delete(tasks []Task, id int) ([]Task, error) {
	if len(tasks) > 0 {
		for i, v := range tasks {
			if v.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				return tasks, nil
			}
		}

	}
	return tasks, errors.New("Нет задач на удаление!")
}
