package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Задача не найдена: %v", err)
			writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			log.Printf("Ошабка базы данных: %v", err)
			writeJSON(w, map[string]string{"error": "Ошибка базы данных"}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}
