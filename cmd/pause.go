/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/HRhades/tk/pkg/database"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause a timer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB(dbPath)
		var timerName string
		if len(args) > 0 {
			timerName = args[0]
		} else {
			log.Fatal("Please supply a name to lookup")
		}

		err := database.Stoptimer(timerName)
		if err != nil {
			log.Fatalf("Stopping timer has failed with error: %v", err)
		} else {
			fmt.Printf("Stopped timer %q", timerName)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
