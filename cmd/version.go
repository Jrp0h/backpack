package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "1.4"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("backpack version %s %s/%s", version, runtime.GOOS, runtime.GOARCH)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}