package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Отсутствует ID задачи"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновляем в БД
	if err := db.UpdateTask(&task); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка обновления задачи"})
		}
		return
	}

	// Успешно — возвращаем пустой объект
	writeJSON(w, map[string]interface{}{})
}
