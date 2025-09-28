package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат JSON"}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Отсутствует ID задачи"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновляем в БД
	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Задача не найдена: %v", err)
			writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			log.Printf("Ошибка обновления задачи: %v", err)
			writeJSON(w, map[string]string{"error": "Ошибка обновления задачи"}, http.StatusInternalServerError)
		}
		return
	}

	// Успешно — возвращаем пустой объект
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
