package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Задача не найдена: %v", err)
			writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			log.Printf("Ошибка удаления задачи: %v", err)
			writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
