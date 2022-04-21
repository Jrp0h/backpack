package cmd

import (
	"fmt"

	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var cryptCmd = &cobra.Command{
	Use:   "crypt",
	Short: "Encrypt/Decrypt files",
	Args: cobra.ExactArgs(1),
	Run: execWithConfig(func (cmd *cobra.Command, args []string, cfg *config.Config){
			
	}),
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt files",
	Args: cobra.ArbitraryArgs,
	Run: encdec(true),
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt files",
	Args: cobra.ArbitraryArgs,
	Run: encdec(false),
}

func encdec(encrypt bool) func (cmd *cobra.Command, args []string) {
	return execWithConfig(func (cmd *cobra.Command, args []string, cfg *config.Config){
		utils.AbortIf(!cfg.Cryption.Enable, "Encryption data is missing from config and is required for this command to run.")

		var f func(string) error
		var verb string

		if encrypt {
			f = cfg.Cryption.Crypter.Encrypt
			verb = "Encrypting"
		} else {
			f = cfg.Cryption.Crypter.Decrypt
			verb = "Decrypting"
		}

		p, _ := pterm.DefaultProgressbar.WithTotal(len(args)).WithTitle(verb + " files").Start()
		p.RemoveWhenDone = true

		// TODO: Accept globs
		for _, path := range args {
			p.UpdateTitle(fmt.Sprintf("%s '%s'", verb, path))
			if err := f(path); err != nil {
				utils.Log.Error("%s", err.Error())
			} else {
				utils.Log.Success("%s", path)
			}
			p.Increment()
		}
	})
}

func init() {
	cryptCmd.AddCommand(encryptCmd)
	cryptCmd.AddCommand(decryptCmd)

	debugCmd.AddCommand(cryptCmd)
}