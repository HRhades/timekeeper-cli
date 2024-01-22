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
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start your paused or stopped timer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB(dbPath)
		var timerName string
		if len(args) > 0 {
			timerName = args[0]
		} else {
			log.Fatal("Please supply a name to lookup")
		}
		timer, err := database.GetTimer(timerName)
		if err != nil {
			log.Fatalf("Reading timer failed: %v", err)
		}
		if timer.Status == "running" {
			fmt.Printf("timer %q is already running", timerName)
			return
		}

		err = database.Starttimer(timerName)
		if err != nil {
			log.Fatalf("Starting timer %q has failed with error: %v", timerName, err)
		} else {
			fmt.Printf("Started timer %q", timerName)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
