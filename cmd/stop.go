package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Register a stop",
	Long:  "Register a stop timestamp if there is an active start timestamp",
	Run: func(cmd *cobra.Command, argv []string) {
		now := time.Now()
		err := Write(Stop, now.Format(time.RFC3339))
		if err != nil {
			fmt.Println(fmt.Errorf("Could not save stop timestamp: %v", err))
		}
		fmt.Println(now.Format(time.RFC3339))
	},
}
