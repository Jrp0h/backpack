package cmd

import (
	"fmt"

	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	testConnectionsCmd = &cobra.Command{
		Use:   "test-connections",
		Short: "Test connections",
		Aliases: []string{"tc"},

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			actions := cfg.Actions.OnlyOrExcept(only, except)

			if len(actions) == 0 {
				utils.Log.Warning("No actions to run.")
				return
			}

			p, _ := pterm.DefaultProgressbar.WithTotal(len(actions)).WithTitle("Testing connections").Start()
			p.RemoveWhenDone = true

			succeded := 0

			for k, v := range actions {
				p.UpdateTitle(fmt.Sprintf("Trying to connect to '%s'", k))
				if err := v.TestConnection(); err != nil {
					utils.Log.Error("%s: %s", k, err.Error())
				} else {
					succeded++
					utils.Log.Success("%s", k)
				}

				p.Increment()
			}

			println()
			switch {
			case succeded == len(actions):
				utils.Log.Success("All actions connected successfully")
			case succeded == 0:
				utils.Log.Error("All actions failed to connect")
			default:
				utils.Log.Warning("%d/%d actions connected successfully", succeded, len(actions))
			}
		}),
	}
)

func init() {
	rootCmd.AddCommand(testConnectionsCmd)

	testConnectionsCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	testConnectionsCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")
}
