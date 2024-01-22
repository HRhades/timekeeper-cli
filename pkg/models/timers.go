package models

import "time"

type Timer struct {
	Id        int64
	Name      string
	Status    string
	Created   int64
	TimerRows []TimerRow
}

type TimerRow struct {
	TimerName       string
	Timestamp_start int64
	Timestamp_end   int64
	TimerDuration   int64
}

func (t Timer) Duration() int64 {
	var duration int64 = 0

	for _, tr := range t.TimerRows {
		if tr.Timestamp_end == 0 {
			currentTime := time.Now()
			tr.TimerDuration = currentTime.Sub(time.UnixMicro(tr.Timestamp_start)).Microseconds()
		}

		duration += tr.TimerDuration
	}

	return duration
}

func (t Timer) LastRow() TimerRow {
	return t.TimerRows[len(t.TimerRows)-1]

	// var returnRow TimerRow
	// for _, tr := range t.TimerRows {
	// 	if tr.Timestamp_end == 0 {
	// 		return tr
	// 	}
	// 	returnRow = tr
	// }
	// return returnRow
}
