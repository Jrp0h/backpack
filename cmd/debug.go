package cmd

import (
	"os"

	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/Jrp0h/backuper/zip"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Different debug commands for testing features",
	Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
		cfg.Cd()

		file := utils.NewFileData("%Y-%m-%d_%H%M", os.TempDir(), "zip")
		err := zip.Zip(cfg.Path, file.Path)
		if err != nil {
			utils.Log.Error("%s", err.Error())
		} else {
			utils.Log.Success("Yaaaay")
		}
	}),
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
