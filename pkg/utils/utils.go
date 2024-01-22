package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/HRhades/tk/pkg/models"
)

func GetTimerDuration(timer models.Timer) string {
	durationMicro := timer.Duration()
	timerDuration := time.Duration(durationMicro) * time.Microsecond

	hour := int(timerDuration.Seconds() / 3600)
	minute := int(timerDuration.Seconds()/60) % 60
	second := int(timerDuration.Seconds()) % 60

	return fmt.Sprintf("%d:%02d:%02d", hour, minute, second)

}

func PrintTimer(timer models.Timer) {
	w := tabwriter.NewWriter(os.Stdout, 3, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tStatus\tDuration (HH:MM:SS)")
	fmt.Fprintln(w, timer.Name+"\t"+timer.Status+"\t"+GetTimerDuration(timer))
	w.Flush()

}

func PrintTimers(timers []models.Timer) {
	w := tabwriter.NewWriter(os.Stdout, 3, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tStatus\tDuration (HH:MM:SS)")
	for _, timer := range timers {
		fmt.Fprintln(w, timer.Name+"\t"+timer.Status+"\t"+GetTimerDuration(timer))
	}
	w.Flush()
}
