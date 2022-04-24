package cmd

import (
	"os"

	"github.com/Jrp0h/backpack/action"
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/handlers"
	"github.com/Jrp0h/backpack/utils"
	"github.com/Jrp0h/backpack/zip"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	rNoBackup bool

	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore from uploaded file",
		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			cfg.Require(config.Actions)
			cfg.Cd()

			if !rNoBackup {
				res, err := handlers.HandleBackup(cfg, handlers.BackupFlags{
					Only:      only,
					Except:    except,
					Force:     force,
					NoEncrypt: noEncrypt,
				})
				if err != nil {
					utils.Log.Fatal("Backup failed. %s", err.Error())
				}

				if res == handlers.BACKUP_NO_ACTIONS {
					utils.Log.Fatal("No actions found. Stopping")
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
				files, err := action.ListFiles()
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

	if !cfg.Crypto.Enabled {
		utils.Log.Info("Encryption Settings isn't set. Skipping decryption.")
		return nil
	}

	return cfg.Crypto.Crypter.Decrypt(path)
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVar(&rNoBackup, "no-backup", false, "Doesn't create backup")
	restoreCmd.Flags().BoolVar(&noEncrypt, "no-encrypt", false, "Doesn't encrypt files")
	restoreCmd.Flags().BoolVar(&force, "force", false, "Force backup even if prev_hash is the same")

	restoreCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	restoreCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")

	restoreCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Path to config file.")
	utils.IgnoreError(restoreCmd.MarkFlagRequired("config"))
}
