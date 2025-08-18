package api

import (
	"database/sql"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"})
		}
		return
	}

	writeJSON(w, map[string]interface{}{})
}
