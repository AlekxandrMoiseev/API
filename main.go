package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task представляет задачу с ее идентификатором, описанием, заметкой и связанными приложениями
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// tasks - это мапа, в которой хранятся наши задачи по их идентификаторам
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks - функция для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	// Установка статуса ответа - OK
	w.WriteHeader(http.StatusOK)

	// Кодирование мапы задач в JSON и отправка в ответ
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		// В случае ошибки кодируем и отправляем статус Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// createTask - функция для создания новой задачи
func createTask(w http.ResponseWriter, r *http.Request) {
	// Декодирование JSON запроса в структуру Task
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// В случае ошибки отправляем статус Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Добавление новой задачи в мапу по ее идентификатору
	tasks[task.ID] = task

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	// Установка статуса ответа - Created
	w.WriteHeader(http.StatusCreated)

	// Кодирование созданной задачи в JSON и отправка в ответ
	if err := json.NewEncoder(w).Encode(task); err != nil {
		// В случае ошибки кодируем и отправляем статус Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getTaskByID - функция для получения задачи по ее идентификатору
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора задачи из URL параметра
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]

	// Если задача с таким идентификатором не найдена, отправляем статус Not Found
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	// Установка статуса ответа - OK
	w.WriteHeader(http.StatusOK)

	// Кодирование найденной задачи в JSON и отправка в ответ
	if err := json.NewEncoder(w).Encode(task); err != nil {
		// В случае ошибки кодируем и отправляем статус Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// deleteTaskByID - функция для удаления задачи по ее идентификатору
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора задачи из URL параметра
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]

	// Если задача с таким идентификатором не найдена, отправляем статус Not Found
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Удаление задачи из мапы по ее идентификатору
	delete(tasks, id)

	// Отправка статуса OK в ответ на успешное удаление
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Инициализация роутера Chi
	r := chi.NewRouter()

	// Регистрация хендлеров для различных маршрутов
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	// Запуск HTTP сервера на порту 8080
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
