package handlers

import (
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

type encryptStatus int

const (
	ENCRYPT_SKIP encryptStatus = iota
	ENCRYPT_CONTINUE
	ENCRYPT_USER_CANCEL
)

func HandleEncrypt(cfg *config.Config, dataPath string, noEncrypt bool) (encryptStatus, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Encrypt: Starting")

	if noEncrypt {
		spinner.Warning("Encrypt: No Encrypt set. Skipping encryption.")
		return ENCRYPT_SKIP, nil
	}

	if !cfg.Crypto.Enabled {
		p := promptui.Prompt{
			Label:     "Encryption isn't set. Are you sure you want to continue?",
			IsConfirm: true,
		}

		_, err := p.Run()

		if err != nil {
			spinner.Warning("Encrypt: User canceled because encryption wasn't set. Stopping.")
			return ENCRYPT_USER_CANCEL, nil
		}

		spinner.Warning("Encrypt: Encryption Settings isn't set. Skipping encryption.")
		return ENCRYPT_SKIP, nil
	}

	err := cfg.Crypto.Crypter.Encrypt(dataPath)
	if err != nil {
		spinner.Fail("Encrypt: Failed to encrypt data")
	}

	return ENCRYPT_CONTINUE, err
}

func HandleDecrypt(cfg *config.Config, dataPath string, noEncrypt bool) (encryptStatus, error) {
	if noEncrypt {
		return ENCRYPT_SKIP, nil
	}

	if !cfg.Crypto.Enabled {
		utils.Log.Info("Encryption Settings isn't set. Skipping decryption.")
		return ENCRYPT_SKIP, nil
	}

	return ENCRYPT_CONTINUE, cfg.Crypto.Crypter.Decrypt(dataPath)

}
