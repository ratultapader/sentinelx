package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {

	var err error

	DB, err = sql.Open("sqlite3", "./sentinelx.db")
	if err != nil {
		return err
	}

	fmt.Println("✅ Database connected")

	return createTables()
}

func createTables() error {

	eventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		type TEXT,
		source_ip TEXT,
		message TEXT
	);`

	alertTable := `
	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		type TEXT,
		severity TEXT,
		source_ip TEXT,
		description TEXT
	);`

	_, err := DB.Exec(eventTable)
	if err != nil {
		return err
	}

	_, err = DB.Exec(alertTable)
	return err
}