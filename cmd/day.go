package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dayCmd)
}

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Display the hour:minutes sum of all the timestamps of the current day",
	Long: `It does the sum of all the timestamp for the current day or displays 00:00 if there are no timestamp.
Active timestamp (start without stop) are also taken into account.`,
	Run: func(cmd *cobra.Command, argv []string) {
		fmt.Println("day")
	},
}
