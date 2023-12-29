package cmd

import (
	"fmt"
	"log"

	"github.com/HRhades/tk/pkg/database"
	"github.com/HRhades/tk/pkg/models"
	"github.com/spf13/cobra"
)

func printTimers(timers []models.TimerRow) {

}

// newCmd represents the new command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all timers",
	Long:  `Use list to list all active timers`,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB("E:\\Coding\\temp_data\\tk.db")
		rowsArray, err := database.GetTimers("all")
		if err != nil {
			log.Fatalf("query failed: %v", err)
		}

		for _, row := range rowsArray {
			fmt.Println(row)
		}

	},
}

var allTimers bool
var pausedTimers bool
var stoppedTimers bool

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allTimers, "allTimers", "a", false, "alltimers:True,False")
	listCmd.Flags().BoolVarP(&pausedTimers, "pausedTimers", "p", false, "pausedTimers:True,False")
	listCmd.Flags().BoolVarP(&stoppedTimers, "stoppedTimers", "s", false, "stoppedTimers:True,False")
}
