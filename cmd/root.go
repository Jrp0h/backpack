package cmd

import (
	"os"

	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/spf13/cobra"
)

var (
	// Shared flags
	only      []string
	except    []string
	noEncrypt bool
	force     bool

	// Global
	cfgPath     string
	debugMode   bool
	verboseMode bool

	RootCmd = &cobra.Command{
		Use:               "backpack",
		Short:             "Easily backup and restore folders to and from different storages",
		DisableAutoGenTag: true,
	}
)

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode. MAY PRINT SENSITIVE INFORMATION")
	RootCmd.PersistentFlags().BoolVar(&verboseMode, "verbose", false, "Print more information.")

	RootCmd.Version = version
}

func execWithConfig(f func(cmd *cobra.Command, args []string, cfg *config.Config)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		utils.Log.DebugEnabled = debugMode
		utils.Log.VerboseEnabled = verboseMode

		cfg, err := config.LoadConfig(cfgPath)
		utils.AbortIfError(err)

		f(cmd, args, cfg)
	}
}
