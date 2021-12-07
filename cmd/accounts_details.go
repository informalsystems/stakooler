package cmd

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/querier"
	"github.com/informalsystems/stakooler/client/display"
	"github.com/informalsystems/stakooler/config"
	"github.com/spf13/cobra"
	"os"
)

// represents the 'accounts details' command
var accountDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Shows detailed information about accounts",
	Long: `This command shows detailed information about configured accounts. For example:

It shows tokens balance, rewards, delegation and unbonding values per account`,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := config.LoadConfig()
		if err != nil {
			fmt.Println("errors reading configuration", err)
			os.Exit(1)
		}
		// Load each account details
		for _, acct := range accounts.Entries {
			acctDetails, err := querier.LoadAccountDetails(acct)
			if err != nil {
				fmt.Println("failed to retrieved",acct.Address,"details:",err)
				os.Exit(1)
			} else {
				acct.AccountDetails = acctDetails
			}

		}
		// Print table information
		display.PrintAccountDetailsTable(&accounts)
	},
}

func init() {
	accountsCmd.AddCommand(accountDetailsCmd)
}
