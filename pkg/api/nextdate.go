package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		var err error
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			http.Error(w, "неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "отсутствует параметр date", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "отсутствует параметр repeat", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", nextDate)
}

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты dstart: %v", err)
	}

	now = TruncateTime(now)
	start = TruncateTime(start)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("пустое правило повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("некорректный формат правила d: требуется число")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("не число в правиле d")
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("интервал должен быть от 1 до 400")
		}

		date := start
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("правило y не должно иметь аргументов")
		}

		date := start
		for {
			date = date.AddDate(1, 0, 0)

			// Коррекция: 29 февраля → 1 марта
			if start.Month() == time.February && start.Day() == 29 {
				if date.Month() == time.February && date.Day() == 28 {
					date = date.AddDate(0, 0, 1)
				}
			}

			if date.After(now) {
				break
			}
		}
		return date.Format(DateLayout), nil

	case "w":
		return "", errors.New("w не реальзован")

	case "m":
		return "", errors.New("m не реальзован")

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", parts[0])
	}
}

func TruncateTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
