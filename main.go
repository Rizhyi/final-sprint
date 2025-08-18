package main

import (
	"log"
	"net/http"

	"github.com/rizhyi/final-sprint/pkg/api"
	"github.com/rizhyi/final-sprint/pkg/db"
)

const (
	port   = "7540"
	webDir = "web"
	dbFile = "scheduler.db"
)

func main() {
	// Инициализируем базу данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.DB.Close()

	// Регистрация API-обработчиков
	api.Init()

	// Обслуживаем статические файлы
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
