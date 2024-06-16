package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task представляет задачу с ее уникальным идентификатором, описанием, заметкой и списком связанных приложений.
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// tasks - это мапа, в которой хранятся наши задачи. Ключом мапы является уникальный идентификатор задачи, а значением - сама структура Task.
var tasks = map[string]Task{}

// getTasks - функция для получения списка всех задач.
// Она кодирует мапу задач в формат JSON для отправки в ответ.
func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// createTask - функция для создания новой задачи.
// Она декодирует JSON запрос в структуру Task, добавляет новую задачу в мапу,
// кодирует созданную задачу в формат JSON для отправки в ответ.
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли уже задача с таким идентификатором
	if _, ok := tasks[task.ID]; ok {
		http.Error(w, "Task with this ID already exists", http.StatusConflict)
		return
	}

	tasks[task.ID] = task

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getTaskByID - функция для получения задачи по ее уникальному идентификатору.
// Она получает идентификатор задачи из URL параметра, ищет задачу в мапе,
// кодирует найденную задачу в формат JSON для отправки в ответ.
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// deleteTaskByID - функция для удаления задачи по ее уникальному идентификатору.
// Она получает идентификатор задачи из URL параметра, ищет задачу в мапе, удаляет ее из мапы,
// и отправляет статус OK в ответ на успешное удаление.
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	delete(tasks, id)
}

func main() {
	// Инициализация роутера Chi для обработки HTTP запросов.
	r := chi.NewRouter()

	// Регистрация хендлеров для различных маршрутов:
	// - GET /tasks - получение списка всех задач
	// - POST /tasks - создание новой задачи
	// - GET /tasks/{id} - получение задачи по идентификатору
	// - DELETE /tasks/{id} - удаление задачи по идентификатору
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	// Запуск HTTP сервера на порту 8080.
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
