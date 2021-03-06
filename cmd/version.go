package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "1.0-dev.1"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Show version. backpack version %s %s/%s", version, runtime.GOOS, runtime.GOARCH),

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("backpack version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}
