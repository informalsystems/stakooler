package cmd

import (
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Displays information about configured accounts",
	Long:  `Displays information about configured accounts such as balances, rewards, staked, unbonding tokens, etc in a table friendly or csv format`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}
