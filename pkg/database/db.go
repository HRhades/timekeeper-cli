package database

import (
	"database/sql"
	"fmt"
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

	_, err = Db.Exec(
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
		rowid           int64
		timerid         int64
		timername       string
		timestamp_start int64
		timestamp_end   sql.NullInt64
		duration        sql.NullInt64
	)
	for rows.Next() {
		err := rows.Scan(&rowid, &timerid, &timername, &timestamp_start, &timestamp_end, &duration)
		if err != nil {
			log.Fatalf("This failed!, %v", err)
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
		Id:        id,
		Name:      name,
		Status:    status,
		TimerRows: timersArray,
	}
	return newTimerResult, nil
}

func GetTimers(filterValue string) ([]models.Timer, error) {
	var sqlBaseString string
	var sqlQueryString string
	var sqlLimitString string

	sqlLimitString = "ORDER BY created DESC LIMIT 10"

	switch filterValue {
	case "all":
		sqlBaseString = "SELECT name FROM timers"
		sqlLimitString = ""
	case "stopped":
		sqlBaseString = "SELECT name FROM timers WHERE status = 'stopped'"
	case "running":
		sqlBaseString = "SELECT name FROM timers WHERE status = 'running'"
	case "paused":
		sqlBaseString = "SELECT name FROM timers WHERE status = 'paused'"
	}
	sqlQueryString = fmt.Sprintf("%s %s", sqlBaseString, sqlLimitString)
	var timersArray []models.Timer

	rows, err := Db.Query(sqlQueryString)
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
	_, err = Db.Exec(`DELETE FROM timer_rows WHERE name = ?`, timerName)
	if err != nil {
		return err
	}
	return nil
}

func Stoptimer(timerName string) error {

	timer, err := GetTimer(timerName)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	timerRow := timer.LastRow()
	if !(timerRow.Timestamp_end == 0) {
		fmt.Printf("Timer %s is already stopped", timerName)
		return nil
	}

	timerDuration := currentTime.Sub(time.UnixMicro(timerRow.Timestamp_start))
	_, err = Db.Exec(`UPDATE timer_rows SET timestamp_end=?, timer_duration=? WHERE timername=? AND timestamp_end IS NULL;`, currentTime.UnixMicro(), timerDuration.Microseconds(), timerName)
	if err != nil {
		return err
	}
	_, err = Db.Exec(`UPDATE timers SET status='stopped' WHERE name=?;`, timerName)
	if err != nil {
		return err
	}

	return nil
}

func Starttimer(timerName string) error {

	currentTime := time.Now().UnixMicro()
	timer, err := GetTimer(timerName)

	if err != nil {
		return err
	}

	if timer.Status == "running" {
		fmt.Printf("timer %q is already running", timerName)
		return nil
	}

	_, err = Db.Exec(`UPDATE timers SET status='running' WHERE name=?;`, timerName)
	if err != nil {
		return err
	}

	_, err = Db.Exec(
		`INSERT INTO timer_rows (timer_id, timername, timestamp_start) VALUES (?,?,?)`, timer.Id, timer.Name, currentTime,
	)
	if err != nil {
		return err
	}

	return nil
}

func Pausetimer(timerName string) error {
	timer, err := GetTimer(timerName)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	timerRow := timer.LastRow()
	if !(timerRow.Timestamp_end == 0) {
		fmt.Printf("Timer %s is already paused", timerName)
		return nil
	}

	timerDuration := currentTime.Sub(time.UnixMicro(timerRow.Timestamp_start))
	_, err = Db.Exec(`UPDATE timer_rows SET timestamp_end=?, timer_duration=? WHERE timername=? AND timestamp_end IS NULL;`, currentTime.UnixMicro(), timerDuration.Microseconds(), timerName)
	if err != nil {
		return err
	}
	_, err = Db.Exec(`UPDATE timers SET status='paused' WHERE name=?;`, timerName)
	if err != nil {
		return err
	}

	return nil
}
