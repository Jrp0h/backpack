package cmd

import (
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Different debug commands for testing features",
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
