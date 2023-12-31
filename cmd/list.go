package cmd

import (
	"log"

	"github.com/HRhades/tk/pkg/database"
	"github.com/HRhades/tk/pkg/utils"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all timers",
	Long:  `Use list to list all active timers`,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB("E:\\Coding\\temp_data\\tk.db")

		var filterString string
		if allTimers {
			filterString = "all"
		} else if stoppedTimers {
			filterString = "stopped"
		} else {
			filterString = "running"
		}

		rowsArray, err := database.GetTimers(filterString)
		if err != nil {
			log.Fatalf("query failed: %v", err)
		}

		utils.PrintTimers(rowsArray)

	},
}

// var pausedTimers bool
var allTimers bool
var stoppedTimers bool

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allTimers, "allTimers", "a", false, "alltimers:True,False")
	// listCmd.Flags().BoolVarP(&pausedTimers, "pausedTimers", "p", false, "pausedTimers:True,False")
	listCmd.Flags().BoolVarP(&stoppedTimers, "stoppedTimers", "s", false, "stoppedTimers:True,False")
}
