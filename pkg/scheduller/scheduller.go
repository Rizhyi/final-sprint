package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи.
// now — текущее время (от него определяется, когда задача должна быть следующей).
// dstart — начальная дата задачи в формате "20060102".
// repeat — правило повторения.
// Возвращает дату в формате "20060102" или ошибку.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	// 1. Проверка на пустое правило
	if repeat == "" {
		return "", errors.New("правило повторения пустое")
	}

	// 2. Парсим начальную дату dstart
	start, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты dstart: %v", err)
	}

	// 3. Разбиваем правило на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("некорректный формат правила: пустое значение")
	}

	// 4. Проверяем поддерживаемые правила: d и y
	switch parts[0] {
	case "d":
		// Проверка: должно быть 2 части: "d" и число
		if len(parts) != 2 {
			return "", errors.New("некорректный формат правила d: требуется число")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("некорректное число дней в правиле d")
		}
		if interval <= 0 {
			return "", errors.New("интервал должен быть положительным")
		}
		if interval > 400 {
			return "", errors.New("максимальный интервал — 400 дней")
		}

		// Сдвигаем дату с шагом interval дней, пока не станет больше now
		date := start
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format("20060102"), nil

	case "y":
		// Проверка: только "y", без аргументов
		if len(parts) != 1 {
			return "", errors.New("некорректный формат правила y: не должно быть аргументов")
		}

		// Сдвигаем дату на год, пока не станет больше now
		date := start
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0) // +1 год
		}
		return date.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", parts[0])
	}
}

// afterNow возвращает true, если date > now (сравнение только по дате, без времени)
func afterNow(date, now time.Time) bool {
	// Обнуляем время для корректного сравнения дат
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly) || dateOnly.Equal(nowOnly)
}
