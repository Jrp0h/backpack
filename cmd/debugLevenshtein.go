package cmd

import (
	"github.com/Jrp0h/backuper/config"
	"github.com/Jrp0h/backuper/utils"
	"github.com/spf13/cobra"
)

var (
	ldCaseInsensitive = false

	levenshteinCmd = &cobra.Command{
		Use:   "levenshtein",
		Short: "Check levenshtein distance between one string and an array of strings",
		Aliases: []string{"ld"},
		Args: cobra.MinimumNArgs(2),
		Run: execWithConfig(func (cmd *cobra.Command, args []string, _ *config.Config){
			needle := args[0]
			haystack := args[1:]

			result := utils.Levenshtein(needle, haystack, ldCaseInsensitive)
			if len(result) == 0 {
				utils.Log.Error("No results found")
				return
			}

			utils.Log.Success("%s are best with an equal distance of %d", utils.JoinSliceAsSentanceStatement(result, func(t utils.LevenshteinResult) string { return t.Text }), result[0].Distance)
		}),
	}
)


func init() {
	levenshteinCmd.Flags().BoolVarP(&ldCaseInsensitive, "ignore-case", "I", false, "Make search case insensitive")
	debugCmd.AddCommand(levenshteinCmd)
}