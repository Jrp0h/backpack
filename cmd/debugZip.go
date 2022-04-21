package cmd

import (
	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/Jrp0h/backuper/zip"
	"github.com/spf13/cobra"
)


var (
	unzip = false

	zipCmd = &cobra.Command{
		Use:   "zip",
		Short: "Zip file or folder",
		Aliases: []string{"ld"},
		Args: cobra.ExactArgs(2),
		Run: execWithConfig(func (cmd *cobra.Command, args []string, _ *config.Config){
			input := args[0]
			output := args[1]

			if unzip {
				if err := zip.Unzip(input, output); err != nil {
					utils.Log.Error("Failed to unzip %s\n%s", input, err)
					return
				}
				utils.Log.Success("%s has been successfully unzipped to %s", input, output)
			} else {
				if err := zip.Zip(input, output); err != nil {
					utils.Log.Error("Failed to zip %s\n%s", input, err)
					return
				}

				utils.Log.Success("%s has been successfully zipped to %s", input, output)
			}
		}),
	}
)


func init() {
	zipCmd.Flags().BoolVarP(&unzip, "unzip", "U", false, "Unzip")
	debugCmd.AddCommand(zipCmd)
}