package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func New(driver, dsn string) (*Storage, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("storage.sqlite.New: %w", err)
	}

	fmt.Println("!")

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date INTEGER NOT NULL,
    	title TEXT NOT NULL,
    	comment TEXT,
    	repeat TEXT CHECK(LENGTH(repeat) <= 128)
	);

	CREATE INDEX IF NOR EXISTS idx_scheduler_date ON scheduler(date);
	`)

	fmt.Println("!")
	if err != nil {
		return nil, fmt.Errorf("storage.sqlite.New: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("storage.sqlite.New: %w", err)
	}

	return &Storage{DB: db}, nil
}
