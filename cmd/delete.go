package cmd

import (
	"fmt"
	"log"

	"github.com/HRhades/tk/pkg/database"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		database.InitDB("E:\\Coding\\temp_data\\tk.db")

		var timerName string
		if len(args) > 0 {
			timerName = args[0]
		}

		err := database.DeleteTimer(timerName)
		if err != nil {
			log.Fatalf("Deletion failed with error: %v", err)
		} else {
			fmt.Printf("Deleted timer: %s", timerName)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
