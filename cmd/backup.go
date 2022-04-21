package cmd

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Jrp0h/backuper/action"
	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/Jrp0h/backuper/zip"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	bNoEncrypt = false
	bForce = false

	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Zip, Encrypt and Store files",

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			cfg.Require(config.Path | config.Actions)
			cfg.Validate(config.Path)
			actions := cfg.Actions.OnlyOrExcept(only, except)

			// Change Current Directory if CWD isn't empty
			if cfg.CWD != "" {
				if err := os.Chdir(cfg.CWD); err != nil {
					utils.Log.Fatal("%s", err.Error());
				}
			}

			file := utils.NewFileData("%Y-%m-%d_%H%M", os.TempDir(), "zip")

			// Zip
			utils.AbortIfError(zip.Zip(cfg.Path, file.Path))
			defer os.Remove(file.Path) // Clean up

			// Hash
			err, shouldContinue := handleHash(cfg, file.Path) 
			if err != nil {
				utils.Log.FatalNoExit(err.Error())
				return
			}

			if !shouldContinue {
				return
			}

			// Encrypt
			err = handleEncrypt(cfg, file.Path)
			if err != nil {
				utils.Log.FatalNoExit(err.Error())
				return
			}

			// Run Actions
			if len(actions) == 0 {
				utils.Log.Warning("No actions to run. Stopping")
				return
			}

			wg := new(sync.WaitGroup)
			m := new(sync.Mutex)
			succeded := 0

			p, _ := pterm.DefaultProgressbar.WithTotal(len(actions)).WithTitle("Running actions").Start()
			p.RemoveWhenDone = true

			for k, v := range actions {
				wg.Add(1)
				go runAction(&file, k, v, p, &succeded, wg, m)
			}

			wg.Wait()

			println()
			switch {
			case succeded == len(actions):
				utils.Log.Success("All actions completed successfully")
			case succeded == 0:
				utils.Log.Error("All actions failed")
			default:
				utils.Log.Warning("%d/%d actions succeded", succeded, len(actions))
			}
		}),
	}
)

func init() {
	backupCmd.Flags().BoolVar(&bNoEncrypt, "no-encrypt", false, "Doesn't encrypt files")
	backupCmd.Flags().BoolVar(&bForce, "force", false, "Force backup even if prev_hash is the same")

	backupCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	backupCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")

	rootCmd.AddCommand(backupCmd)
}

func handleHash(cfg *config.Config, path string) (error, bool) {

	// Compute new hash
	zipped, err := ioutil.ReadFile(path)
	if err != nil {
		return err, false
	}

	h := sha512.New()
	_, err = h.Write(zipped)
	if err != nil {
		return err, false
	}
	newHash := h.Sum(nil)

	// Check Prev Hash
	var prevHash []byte

	if !bForce {
		prev, err := ioutil.ReadFile(cfg.Hash)
		if err != nil {
			return err, false
		}

		prevHash, err = hex.DecodeString(string(prev))
		if err != nil {
			return err, false
		}
	}

	if bytes.Equal(newHash, prevHash) && !bForce {
		utils.Log.Info("Data hasn't changed and force is not enabled. Stopping")
		return nil, false
	}

	// Store new hash
	if utils.PathIsFile(cfg.Hash) {
		err = ioutil.WriteFile(cfg.Hash, []byte(hex.EncodeToString(newHash)), 0644)
		if err != nil {
			return err, false
		}
	}

	return nil, true
}

func handleEncrypt(cfg *config.Config, path string) error {
	if bNoEncrypt {
		return nil
	}

	if !cfg.Cryption.Enabled {
		// TODO: Wait for user input
		utils.Log.Info("Encryption Settings isn't set. Skipping encryption.")
		return nil
	}

	return cfg.Cryption.Crypter.Encrypt(path)
}

func runAction(file *utils.FileData, actionName string, action action.Action, p *pterm.ProgressbarPrinter, succeded *int, wg *sync.WaitGroup, m *sync.Mutex) {
	err := action.Run(file)

	m.Lock()
	if err != nil {
		utils.Log.Error("%s failed\n%s", actionName, err.Error())
	} else {
		*succeded++
		utils.Log.Success("%s succeded", actionName)
	}
	p.Increment()
	m.Unlock()

	wg.Done()
}