package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rizhyi/final-sprint/pkg/db"
)

// writeJSON отправляет JSON-ответ
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// addTaskHandler обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат JSON"})
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка сохранения задачи"})
		return
	}

	// Возвращаем ID
	writeJSON(w, map[string]int64{"id": id})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := TruncateTime(time.Now())

	// Если дата не указана или указан today — ставим сегодня
	if task.Date == "" || strings.ToLower(task.Date) == "today" {
		task.Date = now.Format("20060102")
		return nil
	}

	// Парсим дату
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("некорректный формат даты")
	}

	t = TruncateTime(t)

	// Если дата в прошлом и нет правила повторения — ставим сегодня
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			// Если правило есть — вычисляем следующую дату
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}
