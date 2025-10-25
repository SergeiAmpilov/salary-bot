// internal/storage/storage.go
package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dbPath string) *Storage {
	// Используем драйвер "sqlite" (без "3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Не удалось открыть базу данных:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	storage := &Storage{DB: db}
	storage.initTables()
	return storage
}

func (s *Storage) initTables() {
	query := `
CREATE TABLE IF NOT EXISTS salaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tech TEXT NOT NULL,
    salary_min INTEGER,
    salary_max INTEGER,
    type TEXT NOT NULL,
    experience_min INTEGER,
    experience_max INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal("Не удалось создать таблицу salaries:", err)
	}
	log.Println("Таблица salaries готова")
}
