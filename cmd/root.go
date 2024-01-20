/*
Copyright © 2023 HRhades
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HRhades/tk/pkg/database"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tk",
	Short: "A cli tool to keep time",
	Long:  `A cli tool to keep time`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var dbPath string
var configPath string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// get and set tool dir
	home, err := homedir.Dir()
	if err != nil {
		log.Println("Unable to detect home directory. Please set data file using --dbpath")
	}
	tkDir := filepath.Join(home, ".tk")
	err = os.MkdirAll(tkDir, os.ModePerm)
	if err != nil {
		log.Fatalf("tk dir cannot be made error occurred: %v", err)
	}

	defaultDbPath := filepath.Join(tkDir, "timers.db")
	defaultConfigPath := filepath.Join(tkDir, "config.json")

	rootCmd.PersistentFlags().StringVar(&dbPath, "dbpath", defaultDbPath, "database location (default is $HOME/.tk/timers.db)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "config location (default is $HOME/.tk/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	database.InitDB(dbPath)
}
