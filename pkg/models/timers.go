package models

import "time"

type Timer struct {
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

func (t Timer) duration() int64 {
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
