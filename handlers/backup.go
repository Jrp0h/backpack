package handlers

import (
	"os"
	"sync"

	"github.com/Jrp0h/backpack/action"
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/Jrp0h/backpack/zip"
	"github.com/pterm/pterm"
)

type backupStatus int

const (
	BACKUP_NO_CHANGE backupStatus = iota
	BACKUP_USER_CANCEL
	BACKUP_NO_ACTIONS
	BACKUP_SUCCESS
)

type BackupFlags struct {
	Only   []string
	Except []string

	Force     bool
	NoEncrypt bool
}

func HandleBackup(cfg *config.Config, flags BackupFlags) (backupStatus, error) {
	utils.Log.Debug("handlers/backup: Entered")
	actions := cfg.Actions.OnlyOrExcept(flags.Only, flags.Except)

	// Change Current Directory if CWD isn't empty
	cfg.Cd()

	file := utils.NewFileData(cfg.FileNameFormat, os.TempDir(), "zip")

	// Zip
	utils.AbortIfError(zip.Zip(cfg.Path, file.Path))
	defer os.Remove(file.Path) // Clean up

	// Hash
	// FIXME: Can't restore because data hasn't changed since last backup
	utils.Log.Debug("handlers/backup: Hash start")
	hash, err := HandleHash(cfg, file.Path, flags.Force)
	if err != nil {
		utils.Log.FatalNoExit(err.Error())
		return 0, err
	}
	// Only store hash if backup success
	saveHash := false
	defer func() { utils.IgnoreError(hash.StoreHash(&saveHash)) }()

	if hash.Result == HASH_NO_CHANGE {
		utils.Log.Debug("handlers/backup: Hash no change")
		return BACKUP_NO_CHANGE, nil
	}
	utils.Log.Debug("handlers/backup: Hash done")

	// Encrypt
	utils.Log.Debug("handlers/backup: Encryption start")
	crypt, err := HandleEncrypt(cfg, file.Path, flags.NoEncrypt)
	if err != nil {
		utils.Log.FatalNoExit(err.Error())
		return 0, err
	}

	if crypt == ENCRYPT_USER_CANCEL {
		return BACKUP_USER_CANCEL, nil
	}
	utils.Log.Debug("handlers/backup: Encryption done")

	// Run Actions
	if len(actions) == 0 {
		utils.Log.Warning("No actions to run. Stopping backup")
		// NOTE: Should hash be saved here?
		return BACKUP_NO_ACTIONS, nil
	}

	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	succeded := 0

	p, _ := pterm.DefaultProgressbar.WithTotal(len(actions)).WithTitle("Running actions").Start()
	p.RemoveWhenDone = true

	for k, v := range actions {
		wg.Add(1)
		go backupRunAction(&file, k, v, p, &succeded, wg, m)
	}

	wg.Wait()

	println()
	switch {
	case succeded == len(actions):
		utils.Log.Success("All actions completed successfully")
	case succeded == 0:
		// NOTE: Should hash be saved here?
		utils.Log.Error("All actions failed")
	default:
		utils.Log.Warning("%d/%d actions succeded", succeded, len(actions))
	}

	saveHash = true
	return BACKUP_SUCCESS, nil
}

func backupRunAction(file *utils.FileData, actionName string, action action.Action, p *pterm.ProgressbarPrinter, succeded *int, wg *sync.WaitGroup, m *sync.Mutex) {
	err := action.Upload(file)

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
