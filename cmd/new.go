package cmd

import (
	"log"
	"os"
	"time"

	"github.com/HRhades/tk/pkg/database"
	"github.com/HRhades/tk/pkg/models"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add new timer",
	Long:  `Use new to add a new timer to keep track off`,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB(dbPath)

		var timerName string
		if len(args) > 0 {
			timerName = args[0]
		}
		_, err := database.GetTimer(timerName)

		if err == nil {
			log.Fatalf("A timer with the name %s already exists, please choose another name", timerName)
		}

		newTimer := models.TimerRow{
			Name:            timerName,
			Status:          "running",
			Timestamp_start: time.Now().UnixMicro(),
		}
		id, err := database.AddTimer(&newTimer)
		if err != nil {
			log.Printf("Adding timer failed!: %v", err)
			os.Exit(1)
		}
		log.Printf("Added timer '%v' with id %v", newTimer.Name, id)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
