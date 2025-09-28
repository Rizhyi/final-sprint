package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная переменная для подключения к БД
var DB *sql.DB

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

// Init инициализирует базу данных: открывает или создает файл и таблицу
func Init(dbFile string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(dbFile)
	fileExists := err == nil

	// Открываем базу данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файл был создан, нужно создать таблицу и индекс
	if !fileExists {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
