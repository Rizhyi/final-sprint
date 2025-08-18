package api

import (
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
		writeJSON(w, map[string]string{"error": "ошибка получения задач"})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
