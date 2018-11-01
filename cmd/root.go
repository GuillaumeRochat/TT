package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "TT",
	Short: "Simple start/stop timer to keep track of active time",
	Long: `TT records a timestamp whenever start and stop are called. It then
can display the time difference between all starts and stop within
specific timeframes`,
}

// Execute the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
