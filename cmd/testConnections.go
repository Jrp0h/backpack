package cmd

import (
	"fmt"

	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	only []string

	testConnectionsCmd = &cobra.Command{
		Use:   "test-connections",
		Short: "Test connections",
		Aliases: []string{"tc"},

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			p, _ := pterm.DefaultProgressbar.WithTotal(len(cfg.Actions)).WithTitle("Testing connections").Start()
			p.RemoveWhenDone = true

			for k, v := range cfg.Actions {
				p.UpdateTitle(fmt.Sprintf("Trying to connect to '%s'", k))
				if err := v.TestConnection(); err != nil {
					utils.Log.Error("%s: %s", k, err.Error())
				} else {
					utils.Log.Success("%s", k)
				}

				p.Increment()
			}
		}),
	}
)

func init() {
	rootCmd.AddCommand(testConnectionsCmd)

	testConnectionsCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
}
