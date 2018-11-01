package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(weekCmd)
}

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Display the hour:minute sum of all the timestamps of the current week",
	Long: `It does the sum of all the timestamps for the current week or displays 00:00 if there are no timestamp.
The week beging on Sunday. Active timestamps (start without stop) are also taken into account.`,
	Run: func(cmd *cobra.Command, argv []string) {
		fmt.Println("week")
	},
}
