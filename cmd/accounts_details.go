package cmd

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/informalsystems/stakooler/client/cosmos/querier"
	"github.com/informalsystems/stakooler/client/display"
	"github.com/informalsystems/stakooler/config"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var flagCsv *bool

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

		var wg sync.WaitGroup

		// Progress bar
		// iterations are the api calls number times the number of accounts
		total_iterations := len(accounts.Entries) * 6
		bar := progressbar.NewOptions(total_iterations,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(25),
			progressbar.OptionUseANSICodes(false),
			progressbar.OptionSetDescription("Getting accounts details..."),
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "▪︎[reset]",
				SaucerHead:    ">[reset]",
				SaucerPadding: ".",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		// Load each account details
		for _, acct := range accounts.Entries {
			wg.Add(1)
			bar.Describe(fmt.Sprintf("Getting account %s details", acct.Address))
			go func(account *model.Account) {
				err := querier.LoadTokenInfo(account)
				defer wg.Done()
				if err != nil {
					fmt.Println("failed to retrieved", acct.Address, "details:", err)
					os.Exit(1)
				}
			}(acct)
		}

		// Wait for load token info tasks to finish
		wg.Wait()

		// If csv flag specified use csv output
		if *flagCsv {
			// write csv file
			display.WriteCSV(&accounts)
		} else {
			// Print table information
			display.PrintAccountDetailsTable(&accounts)
		}
	},
}

func init() {
	flagCsv = accountDetailsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	accountsCmd.AddCommand(accountDetailsCmd)
}
