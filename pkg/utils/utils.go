package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/HRhades/tk/pkg/models"
)

func GetTimerDuration(timer models.TimerRow) string {
	var timerStart time.Time
	var timerEnd time.Time
	var timerDuration time.Duration

	timerStart = time.UnixMicro(timer.Timestamp_start)

	if timer.Timestamp_end == 0 {
		timerEnd = time.Now()
	} else {
		timerEnd = time.UnixMicro(timer.Timestamp_end)
	}

	timerDuration = timerEnd.Sub(timerStart)

	hour := int(timerDuration.Seconds() / 3600)
	minute := int(timerDuration.Seconds()/60) % 60
	second := int(timerDuration.Seconds()) % 60

	return fmt.Sprintf("%d:%02d:%02d", hour, minute, second)

}

func PrintTimer(timer models.TimerRow) {
	w := tabwriter.NewWriter(os.Stdout, 3, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tStatus\tDuration (HH:MM:SS)")
	fmt.Fprintln(w, timer.Name+"\t"+timer.Status+"\t"+GetTimerDuration(timer))
	w.Flush()

}

func PrintTimers(timers []models.TimerRow) {
	w := tabwriter.NewWriter(os.Stdout, 3, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tStatus\tDuration (HH:MM:SS)")
	for _, row := range timers {
		fmt.Fprintln(w, row.Name+"\t"+row.Status+"\t"+GetTimerDuration(row))
	}
	w.Flush()
}
