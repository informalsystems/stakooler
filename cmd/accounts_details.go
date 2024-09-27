package cmd

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/querier"
	"github.com/informalsystems/stakooler/client/display"
	"github.com/informalsystems/stakooler/config"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	flagCsv            *bool
	flagZbxAcctDetails *bool
)

// represents the 'accounts details' command
var accountDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Shows detailed information about accounts",
	Long: `This command shows detailed information about configured accounts. For example:

It shows tokens balance, rewards, delegation and unbonding values per account`,
	Run: func(cmd *cobra.Command, args []string) {
		barEnabled := !*flagCsv && !*flagZbxAcctDetails
		tomlConfig, err := config.LoadConfig(flagConfigPath)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading configuration file:")
		}

		// Load account data
		rawAcctData, err := config.ReadAccountData("")
		if err != nil {
			log.Fatal().Err(err).Msg("error reading account data file")
		}

		var bar *progressbar.ProgressBar

		if barEnabled {
			// Progress bar
			// iterations are the api calls number times the number of accounts
			totalIterations := len(tomlConfig.Accounts.Entries) * 6
			bar = progressbar.NewOptions(totalIterations, progressbar.OptionEnableColorCodes(true), progressbar.OptionShowBytes(false), progressbar.OptionSetWidth(25), progressbar.OptionUseANSICodes(false), progressbar.OptionClearOnFinish(), progressbar.OptionSetPredictTime(false), progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "▪︎[reset]",
				SaucerHead:    ">[reset]",
				SaucerPadding: ".",
				BarStart:      "[",
				BarEnd:        "]",
			}))
		} else {
			bar = progressbar.New(0)
		}

		httpClient := cosmos.NewHttpClient()
		chains := config.ParseChainConfig(rawAcctData, httpClient)

		for _, chain := range chains {
			for _, account := range chain.Accounts {
				if barEnabled {
					bar.Describe(fmt.Sprintf("Getting account %s details", account.Name))
				}
			}
		}

		// Load each account details
		for _, acct := range tomlConfig.Accounts.Entries {

			// Don't show this if csv option enabled
			if barEnabled {
				bar.Describe(fmt.Sprintf("Getting account %s details", acct.Address))
			}

			// Get latest block information to include in the account
			blockInfo, err := api.GetLatestBlock(acct.Chain, httpClient)
			if err != nil {
				bar.Describe(fmt.Sprintf("failed to get latest block: %s", err))
			}
			bar.Add(1)

			acct.BlockHeight = blockInfo.Block.Header.Height
			acct.BlockTime = blockInfo.Block.Header.Time

			err = querier.LoadAuthData(acct, httpClient)
			if err != nil {
				bar.Describe(err.Error())
			}

			err = querier.LoadBankBalances(acct, httpClient)
			if err != nil {
				bar.Describe(err.Error())
			}
			bar.Add(1)

			err = querier.LoadDistributionData(acct, httpClient)
			if err != nil {
				bar.Describe(err.Error())
			}
			bar.Add(1)

			err = querier.LoadStakingData(acct, httpClient)
			if err != nil {
				bar.Describe(err.Error())
			}
			bar.Add(1)
		}

		// Hide bar
		bar.Finish()

		// If csv flag specified use csv output
		if *flagCsv {
			// write csv file
			display.WriteAccountsCSV(&tomlConfig.Accounts)
		} else if *flagZbxAcctDetails {
			display.ZbxSendChainDiscovery(&tomlConfig)
			display.ZbxSendAccountsDiscovery(&tomlConfig)
			display.ZbxAccountsDetails(&tomlConfig)
		} else {
			// Print table information
			display.PrintAccountDetailsTable(&tomlConfig.Accounts)
		}
	},
}

func init() {
	flagCsv = accountDetailsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	flagZbxAcctDetails = accountDetailsCmd.Flags().BoolP("zbx", "z", false, "push the result to a zabbix trapper item")
	accountsCmd.AddCommand(accountDetailsCmd)
}
