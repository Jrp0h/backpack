/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/Jrp0h/backpack/action"
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/Jrp0h/backpack/zip"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var (
	rNoBackup bool

	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore from uploaded file",
		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			cfg.Require(config.Actions)
			cfg.Cd()

			if !rNoBackup {
				if !backup(cfg) {
					utils.Log.Info("Backup stopped. Stopping")
					return
				}
			}


			var action action.Action
			var file string

			for {
				// TODO: Add --action flag
				actionPrompt := promptui.Select{
					Label: "Select from which action the backup should be restored from",
					Items: cfg.Actions.Names(),
				}

				_, result, err := actionPrompt.Run()
				utils.AbortIfError(err)

				action = cfg.Actions[result]	
				files, err := action.ListFiles();
				utils.AbortIfError(err)

				if len(files) == 0 {
					utils.Log.Warning("%s has no files. Please select again", result)
					continue
				}

				// TODO: Add --file flag
				filePrompt := promptui.Select{
					Label: "Select file",
					Items: files,
				}

				_, file, err = filePrompt.Run()
				utils.AbortIfError(err)
				break
			}

			// Fetch File
			fetchedPath, err := action.Fetch(file)
			utils.AbortIfError(err)
			defer os.Remove(fetchedPath)

			// Decrypt
			if err = handleDecrypt(cfg, fetchedPath); err != nil {
				utils.Log.FatalNoExit(err.Error())
				return
			}

			// Remove old
			os.RemoveAll(cfg.Path)

			// Unzip
			if err = zip.Unzip(fetchedPath, cfg.Path); err != nil {
				utils.Log.FatalNoExit(err.Error())
				return
			}

			utils.Log.Success("Data has been restored!")
		}),
	}
)

func handleDecrypt(cfg *config.Config, path string) error {
	if noEncrypt {
		return nil
	}

	if !cfg.Cryption.Enabled {
		utils.Log.Info("Encryption Settings isn't set. Skipping decryption.")
		return nil
	}

	return cfg.Cryption.Crypter.Decrypt(path)
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVar(&rNoBackup, "no-backup", false, "Doesn't create backup")
	restoreCmd.Flags().BoolVar(&noEncrypt, "no-encrypt", false, "Doesn't encrypt files")
	restoreCmd.Flags().BoolVar(&force, "force", false, "Force backup even if prev_hash is the same")

	restoreCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	restoreCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")
}
