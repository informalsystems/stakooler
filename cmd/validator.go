package cmd

import (
	"github.com/spf13/cobra"
)

// validatorCmd represents the validator command
var validatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "Displays information about a validator",
	Long:  `Displays information about a validator such as voting power, rank, number of delegators, etc in a table friendly or csv format`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(validatorCmd)
}
