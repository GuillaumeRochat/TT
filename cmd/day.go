package cmd

import (
	"fmt"
	"time"

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
		now := time.Now()
		begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 99, time.UTC)
		duration, err := Read(begin.Format(time.RFC3339), end.Format(time.RFC3339))
		if err != nil {
			fmt.Println(fmt.Errorf("Could not read timestamps: %v", err))
		}
		fmt.Println(duration)
	},
}
