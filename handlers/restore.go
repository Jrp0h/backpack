package handlers

import (
	"os"

	"github.com/Jrp0h/backpack/action"
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/Jrp0h/backpack/zip"
	"github.com/manifoldco/promptui"
)

type restoreStatus int

const (
	RESTORE_USER_CANCEL restoreStatus = iota
	RESTORE_NO_ACTIONS
	RESTORE_SUCCESS
)

type RestoreFlags struct {
	Only   []string
	Except []string

	Force     bool
	NoEncrypt bool
	NoBackup  bool

	Action string
	File   string
}

func HandleRestore(cfg *config.Config, flags RestoreFlags) (restoreStatus, error) {

	if !flags.NoBackup {
		res, err := HandleBackup(cfg, BackupFlags{
			Only:      flags.Only,
			Except:    flags.Except,
			Force:     flags.Force,
			NoEncrypt: flags.NoEncrypt,
		})
		if err != nil {
			utils.Log.Fatal("Backup failed. %s", err.Error())
		}

		if res == BACKUP_NO_ACTIONS {
			return RESTORE_NO_ACTIONS, nil
		}

		if res == BACKUP_USER_CANCEL {
			return RESTORE_USER_CANCEL, nil
		}
	}

	var action action.Action
	var file string

	action, file, err := restoreGetFileAndAction(cfg, flags.Action, flags.File)
	if err != nil {
		return 0, err
	}

	// Fetch File
	fetchedPath, err := action.Fetch(file)
	if err != nil {
		return 0, err
	}
	defer os.Remove(fetchedPath)

	// Decrypt
	_, err = HandleDecrypt(cfg, fetchedPath, flags.NoEncrypt)
	if err != nil {
		return 0, err
	}

	// Move old to be able to unzip
	err = os.Rename(cfg.Path, cfg.Path+".tmp")
	if err != nil {
		return 0, err
	}

	// Unzip
	if err = zip.Unzip(fetchedPath, cfg.Path); err != nil {
		osErr := os.Rename(cfg.Path+".tmp", cfg.Path) // Move the file back
		if osErr != nil {
			return 0, osErr
		}

		return 0, err
	}

	// remove tmp file
	err = os.RemoveAll(cfg.Path + ".tmp")
	if err != nil {
		return 0, err
	}

	return RESTORE_SUCCESS, nil
}

func restoreGetFileAndAction(cfg *config.Config, actionName string, file string) (action.Action, string, error) {
	var action action.Action = nil
	var files []string

	if actionName != "" {
		a, exists := cfg.Actions[actionName]

		if exists {
			f, err := a.ListFiles()
			if err != nil {
				return nil, "", err
			}
			files = f

			if len(files) == 0 {
				utils.Log.Warning("%s has no files. Please select another action.", actionName)
			} else {
				action = a
			}
		} else {
			alternatives := utils.Levenshtein(actionName, cfg.Actions.Names(), true).AsQuestion()
			utils.Log.Error("No action named %s, did you mean %s?", actionName, alternatives)
		}
	}

	for {
		if action == nil {
			actionPrompt := promptui.Select{
				Label: "Select from which action the backup should be restored from",
				Items: cfg.Actions.Names(),
			}

			_, result, err := actionPrompt.Run()
			if err != nil {
				return nil, "", err
			}

			action := cfg.Actions[result]
			files, err = action.ListFiles()
			if err != nil {
				return nil, "", err
			}
			utils.Log.Debug("%s", files)
			if len(files) == 0 {
				utils.Log.Warning("%s has no files. Please select again", result)
				action = nil
				continue
			}
		}

		if file != "" {
			return action, file, nil
		}

		filePrompt := promptui.Select{
			Label: "Select file",
			Items: files,
		}

		_, file, err := filePrompt.Run()
		if err != nil {
			return nil, "", err
		}

		return action, file, nil
	}
}
