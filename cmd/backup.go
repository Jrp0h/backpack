package cmd

import (
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/handlers"
	"github.com/Jrp0h/backpack/utils"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Zip, Encrypt and Store files",

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			cfg.Require(config.Path | config.Actions)
			_, err := handlers.HandleBackup(cfg, handlers.BackupFlags{
				Only:      only,
				Except:    except,
				Force:     force,
				NoEncrypt: noEncrypt,
			})

			if err != nil {
				utils.Log.Fatal("Backup failed: %s", err.Error())
			}

			utils.Log.Debug("Done")
		}),
	}
)

func init() {
	backupCmd.Flags().BoolVar(&noEncrypt, "no-encrypt", false, "Doesn't encrypt files")
	backupCmd.Flags().BoolVar(&force, "force", false, "Force backup even if prev_hash is the same")

	backupCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	backupCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")

	backupCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Path to config file.")
	utils.IgnoreError(backupCmd.MarkFlagRequired("config"))

	rootCmd.AddCommand(backupCmd)
}
