package cmd

import (
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/handlers"
	"github.com/Jrp0h/backpack/utils"
	"github.com/spf13/cobra"
)

var (
	rNoBackup bool
	rAction   string
	rFile     string

	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore from uploaded file",
		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			cfg.Require(config.Actions)
			cfg.Cd()
			_, err := handlers.HandleRestore(cfg, handlers.RestoreFlags{
				Only:      only,
				Except:    except,
				Force:     force,
				NoEncrypt: noEncrypt,
				NoBackup:  rNoBackup,
				Action:    rAction,
				File:      rFile,
			})
			if err != nil {
				utils.Log.Fatal("Restore failed. %s", err.Error())
			}
		}),
	}
)

func init() {
	RootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVar(&rNoBackup, "no-backup", false, "Doesn't create backup")
	restoreCmd.Flags().BoolVar(&noEncrypt, "no-encrypt", false, "Doesn't encrypt files")
	restoreCmd.Flags().BoolVar(&force, "force", false, "Force backup even if prev_hash is the same")

	restoreCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	restoreCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")

	restoreCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Path to config file.")
	restoreCmd.Flags().StringVarP(&rAction, "action", "a", "", "Name of action to restore from.")
	restoreCmd.Flags().StringVarP(&rFile, "file", "f", "", "Name of file to restore from")
	utils.IgnoreError(restoreCmd.MarkFlagRequired("config"))
}
