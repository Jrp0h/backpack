package cmd

import (
	"fmt"

	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	testConnectionsCmd = &cobra.Command{
		Use:     "test-connections",
		Short:   "Test connections",
		Aliases: []string{"tc"},

		Run: execWithConfig(func(cmd *cobra.Command, args []string, cfg *config.Config) {
			actions := cfg.Actions.OnlyOrExcept(only, except)
			cfg.Cd()

			if len(actions) == 0 {
				utils.Log.Warning("No actions to run.")
				return
			}

			p, _ := pterm.DefaultProgressbar.WithTotal(len(actions)).WithTitle("Testing connections").Start()
			p.RemoveWhenDone = true

			succeded := 0
			unknown := 0

			for k, v := range actions {
				p.UpdateTitle(fmt.Sprintf("Trying to connect to '%s'", k))
				if !v.CanValidateConnection() {
					utils.Log.Info("%s: can't test connection", k)
					unknown++
					p.Increment()
					continue
				}
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
				utils.Log.Info("%d/%d actions connected successfully. %d unknown", succeded, len(actions), unknown)
			}
		}),
	}
)

func init() {
	RootCmd.AddCommand(testConnectionsCmd)

	testConnectionsCmd.Flags().StringArrayVar(&only, "only", []string{}, "List of connections to try.")
	testConnectionsCmd.Flags().StringArrayVar(&except, "except", []string{}, "List of connections to ignore.")

	testConnectionsCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Path to config file.")
	utils.IgnoreError(testConnectionsCmd.MarkFlagRequired("config"))
}
