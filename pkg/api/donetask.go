package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/rizhyi/final-sprint/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// 1. Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка базы данных"})
		}
		return
	}

	now := time.Now()

	// 2. Если нет правила повторения — удаляем
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"})
			return
		}
	} else {
		// 3. Если есть правило — вычисляем следующую дату
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка NextDate: %v", err)
			writeJSON(w, map[string]string{"error": "Некорректное правило повторения"})
			return
		}

		// 4. Обновляем ТОЛЬКО дату
		err = db.UpdateDate(id, nextDate)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Ошибка обновления даты задачи"})
			return
		}
	}

	// 5. Успешно — пустой JSON
	writeJSON(w, map[string]interface{}{})
}
