package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func newSQLiteAdapter(file string) *Adapter {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
	return &Adapter{
		Database: db,
		URI:      file,
		Type:     "sqlite",
	}
}
