package cmd

import (
	"fmt"
	"time"

	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/spf13/cobra"
)

var (
	testFormatCmd = &cobra.Command{
		Use:     "test-format",
		Short:   "Test format",
		Aliases: []string{"tc"},

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			fmt.Println(utils.FormatDate(cfg.FileNameFormat, time.Now()))
		}),
	}
)

func init() {
	RootCmd.AddCommand(testFormatCmd)

	testFormatCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Path to config file.")
	utils.IgnoreError(testFormatCmd.MarkFlagRequired("config"))
}
