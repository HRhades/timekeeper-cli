package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/HRhades/tk/pkg/models"
	_ "modernc.org/sqlite"
)

var Db *sql.DB

func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS timers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            status TEXT NOT NULL,
			timestamp_start INTEGER NOT NULL,
			timestamp_end INTEGER
        )`)

	Db = db

	return nil
}

func AddTimer(tr *models.TimerRow) (int64, error) {
	fmt.Println(tr)
	result, err := Db.Exec(
		`INSERT INTO timers (name, status, timestamp_start) VALUES (?,?,?);`, tr.Name, tr.Status, tr.Timestamp_start,
	)
	if err != nil {
		log.Fatalf("insertion failed, %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
