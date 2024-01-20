package database

import (
	"database/sql"
	"log"
	"time"

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
			timer_time INTEGER
        )`)

	Db = db

	return nil
}

func AddTimer(tr *models.TimerRow) (int64, error) {
	// fmt.Println(tr)
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

func GetTimer(timerName string) (models.TimerRow, error) {
	row := Db.QueryRow(`SELECT * FROM timers WHERE name = ?;`, timerName)
	var (
		id              int64
		name            string
		status          string
		timestamp_start int64
		timestamp_end   sql.NullInt64
	)
	err := row.Scan(&id, &name, &status, &timestamp_start, &timestamp_end)
	// fmt.Println(id, name, status, timestamp_start, timestamp_end)
	if err != nil {
		return models.TimerRow{}, err
	}

	newRowResult := models.TimerRow{
		Name:            name,
		Status:          status,
		Timestamp_start: timestamp_start,
		Timestamp_end:   timestamp_end.Int64,
	}
	return newRowResult, nil
}

func GetTimers(filterValue string) ([]models.TimerRow, error) {
	var sqlString string
	switch filterValue {
	case "all":
		sqlString = "SELECT * FROM timers"
	case "stopped":
		sqlString = "SELECT * FROM timers WHERE status = 'stopped'"
	case "running":
		sqlString = "SELECT * FROM timers WHERE status = 'running'"
	case "paused":
		sqlString = "SELECT * FROM timers WHERE status = 'paused'"
	}
	var timersArray []models.TimerRow

	rows, err := Db.Query(sqlString)
	if err != nil {
		return timersArray, err
	}
	defer rows.Close()
	var (
		id              int64
		name            string
		status          string
		timestamp_start int64
		timestamp_end   sql.NullInt64
	)
	for rows.Next() {
		err := rows.Scan(&id, &name, &status, &timestamp_start, &timestamp_end)
		if err != nil {
			return timersArray, err
		}
		// fmt.Println(id, name, status, timestamp_start, timestamp_end)

		newRowResult := models.TimerRow{
			Name:            name,
			Status:          status,
			Timestamp_start: timestamp_start,
			Timestamp_end:   timestamp_end.Int64,
		}
		timersArray = append(timersArray, newRowResult)

	}
	return timersArray, nil
}

func DeleteTimer(timerName string) error {
	_, err := Db.Exec(`DELETE FROM timers WHERE name = ?`, timerName)

	if err != nil {
		return err
	}
	return nil
}

func Stoptimer(timerName string) error {

	currentTime := time.Now().UnixMicro()
	_, err := Db.Exec(`UPDATE timers SET status='stopped', timestamp_end=? WHERE name=?;`, currentTime, timerName)

	if err != nil {
		return err
	}

	return nil
}

func Pausetimer(timerName string) error {

	timerRow, err := GetTimer(timerName)
	if err != nil {
		return err
	}
	currentTime := time.Now().UnixMicro()
	timerDuration := currentTime.Sub(timerRow.Timestamp_start)

	_, err := Db.Exec(`UPDATE timers SET status='stopped', timestamp_end=? WHERE name=?;`, currentTime, timerName)

	if err != nil {
		return err
	}

	return nil
}
