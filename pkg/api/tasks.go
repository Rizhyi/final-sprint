package api

import (
	"log"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

// TasksResp — структура JSON-ответа
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler — обработчик GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	const maxLimit = 50

	tasks, err := db.Tasks(maxLimit)
	if err != nil {
		log.Printf("Ошибка получения задачи: %v", err)
		writeJSON(w, map[string]string{"error": "ошибка получения задач"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
