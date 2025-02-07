package storage

import (
	"database/sql"
	"fmt"
)

func OpenDB(driver string, dsn string) (*sql.DB, func(), error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open database: %w", err)
	}

	closeFunc := func() {
		_ = db.Close()
	}

	return db, closeFunc, nil
}
