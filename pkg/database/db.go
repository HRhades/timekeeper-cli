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
            status TEXT NOT NULL
            created INTEGER NOT NULL
        )`)

	db.Exec(`CREATE TABLE IF NOT EXISTS timer_rows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timer_id INTEGER,
			timername TEXT NOT NULL,
			timestamp_start INTEGER NOT NULL,
			timestamp_end INTEGER,
			timer_duration INTEGER
        )`)

	Db = db

	return nil
}

func AddTimer(tr models.Timer) (int64, error) {
	// fmt.Println(tr)
	result, err := Db.Exec(
		`INSERT INTO timers (name, status, created) VALUES (?,?,?);`, tr.Name, tr.Status, tr.Created,
	)
	if err != nil {
		log.Fatalf("insertion failed, %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	result, err = Db.Exec(
		`INSERT INTO timer_rows (timer_id, timername, timestamp_start) VALUES (?,?,?)`, id, tr.Name, tr.Created,
	)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTimer(timerName string) (models.Timer, error) {
	row := Db.QueryRow(`SELECT * FROM timers WHERE name = ?;`, timerName)
	var (
		id      int64
		name    string
		status  string
		created int64
	)

	err := row.Scan(&id, &name, &status, &created)
	// fmt.Println(id, name, status, timestamp_start, timestamp_end)
	if err != nil {
		return models.Timer{}, err
	}

	var timersArray []models.TimerRow

	rows, err := Db.Query(`SELECT * FROM timer_rows WHERE timername = ?`, timerName)
	if err != nil {
		return models.Timer{}, err
	}
	defer rows.Close()
	var (
		timerid         int64
		timername       string
		timestamp_start int64
		timestamp_end   sql.NullInt64
		duration        sql.NullInt64
	)
	for rows.Next() {
		err := rows.Scan(&timerid, &timername, &timestamp_start, &timestamp_end, duration)
		if err != nil {
			return models.Timer{}, err
		}
		// fmt.Println(id, name, status, timestamp_start, timestamp_end)

		newRowResult := models.TimerRow{
			TimerName:       timername,
			Timestamp_start: timestamp_start,
			Timestamp_end:   timestamp_end.Int64,
			TimerDuration:   duration.Int64,
		}
		timersArray = append(timersArray, newRowResult)

	}

	newTimerResult := models.Timer{
		Name:      name,
		Status:    status,
		TimerRows: timersArray,
	}
	return newTimerResult, nil
}

func GetTimers(filterValue string) ([]models.Timer, error) {
	var sqlString string
	switch filterValue {
	case "all":
		sqlString = "SELECT name FROM timers"
	case "stopped":
		sqlString = "SELECT name FROM timers WHERE status = 'stopped'"
	case "running":
		sqlString = "SELECT name FROM timers WHERE status = 'running'"
	case "paused":
		sqlString = "SELECT name FROM timers WHERE status = 'paused'"
	}
	var timersArray []models.Timer

	rows, err := Db.Query(sqlString)
	if err != nil {
		return timersArray, err
	}
	defer rows.Close()
	var name string
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			return timersArray, err
		}

		timer, err := GetTimer(name)
		if err != nil {
			return timersArray, err
		}
		timersArray = append(timersArray, timer)

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
	currentTime := time.Now()
	timerDuration := currentTime.Sub(time.UnixMicro(timerRow.Timestamp_start))

	_, err = Db.Exec(`UPDATE timers SET status='paused', timestamp_end=?, timer_duration=? WHERE name=?;`, currentTime, timerDuration.Microseconds(), timerName)

	if err != nil {
		return err
	}

	return nil
}
