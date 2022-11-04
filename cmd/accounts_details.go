package cmd

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/querier"
	"github.com/informalsystems/stakooler/client/display"
	"github.com/informalsystems/stakooler/config"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"os"
)

var flagCsv *bool

// represents the 'accounts details' command
var accountDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Shows detailed information about accounts",
	Long: `This command shows detailed information about configured accounts. For example:

It shows tokens balance, rewards, delegation and unbonding values per account`,
	Run: func(cmd *cobra.Command, args []string) {
		barEnabled := !*flagCsv
		config, err := config.LoadConfig(flagConfigPath)
		if err != nil {
			fmt.Println("error reading configuration file:", err)
			os.Exit(1)
		}

		var bar *progressbar.ProgressBar

		if barEnabled {
			// Progress bar
			// iterations are the api calls number times the number of accounts
			totalIterations := len(config.Accounts.Entries) * 6
			bar = progressbar.NewOptions(totalIterations,
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetWidth(25),
				progressbar.OptionUseANSICodes(false),
				progressbar.OptionClearOnFinish(),
				progressbar.OptionSetPredictTime(false),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "▪︎[reset]",
					SaucerHead:    ">[reset]",
					SaucerPadding: ".",
					BarStart:      "[",
					BarEnd:        "]",
				}))
		} else {
			bar = progressbar.New(0)
		}

		// Load each account details
		for _, acct := range config.Accounts.Entries {

			// Don't show this if csv option enabled
			if barEnabled {
				bar.Describe(fmt.Sprintf("Getting account %s details", acct.Address))
			}

			err := querier.LoadTokenInfo(acct, bar)
			if err != nil {
				bar.Describe(fmt.Sprintf("failed to retrieve %s details: %s", acct.Address, err))
				//os.Exit(1)
			} else {
				// Don't show this if csv option enabled
				if barEnabled {
					bar.Describe(fmt.Sprintf("Got account %s details", acct.Address))
				}
			}
		}

		// Hide bar
		bar.Finish()

		// If csv flag specified use csv output
		if *flagCsv {
			// write csv file
			display.WriteAccountsCSV(&config.Accounts)
		} else {
			// Print table information
			display.PrintAccountDetailsTable(&config.Accounts)
		}
	},
}

func init() {
	flagCsv = accountDetailsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	accountsCmd.AddCommand(accountDetailsCmd)
}
