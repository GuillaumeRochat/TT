package cmd

import (
	"fmt"
	"time"

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
		now := getOffsetNow()
		begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 99, time.UTC).AddDate(0, 0, 7)
		duration, err := Read(begin.Format(time.RFC3339), end.Format(time.RFC3339))
		if err != nil {
			fmt.Println(fmt.Errorf("Could not read timestamps: %v", err))
		}
		fmt.Println(duration)
	},
}

func getOffsetNow() time.Time {
	now := time.Now()
	offset := 0

	switch now.Weekday() {
	case time.Monday:
		offset = 1
	case time.Tuesday:
		offset = 2
	case time.Wednesday:
		offset = 3
	case time.Thursday:
		offset = 4
	case time.Friday:
		offset = 5
	case time.Saturday:
		offset = 6
	case time.Sunday:
		offset = 7
	}

	return now.AddDate(0, 0, -offset)
}
