package cmd

import (
	"fmt"
	"log"
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
		database.InitDB("E:\\Coding\\temp_data\\tk.db")

		for _, arg := range args {
			fmt.Printf("arg value: %v\n", arg)
		}

		newTimer := models.TimerRow{
			Name:            "testTimer",
			Status:          "running",
			Timestamp_start: time.Now().UnixMicro(),
		}
		id, err := database.AddTimer(&newTimer)
		if err != nil {
			log.Printf("Adding timer failed!: %v", err)
		}
		log.Printf("Added timer '%v' with id %v", newTimer.Name, id)

		fmt.Printf("new called at %v", time.Now().Format(time.RFC3339))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
