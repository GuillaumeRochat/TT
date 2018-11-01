package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Register a start",
	Long:  "Register a start timestamp is there is no other active start timestamp",
	Run: func(cmd *cobra.Command, argv []string) {
		now := time.Now()
		err := Write(Start, now.Format(time.RFC3339))
		if err != nil {
			fmt.Println(fmt.Errorf("Could not save start timestamp: %v", err))
		}
		fmt.Println(now.Format(time.RFC3339))
	},
}
