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

		for _, chain := range chains.Entries {
			chain.AssetList, err = api.GetAssetsList(chain.Name, httpClient)
			if err != nil {
				log.Error().Err(err).Msg("error getting asset list")
			}

			blockInfo, err := api.GetLatestBlock(chain.RestEndpoint, httpClient)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("failed to get latest block, skipping chain %s", chain.Id))
			}

			for _, account := range chain.Accounts {
				account.BlockTime = blockInfo.Block.Header.Time
				account.BlockHeight = blockInfo.Block.Header.Height

				if barEnabled {
					bar.Describe(fmt.Sprintf("Getting account %s details", account.Name))
				}

				if err = querier.LoadAuthData(account, httpClient, chain); err != nil {
					bar.Describe(err.Error())
				}
				bar.Add(1)

				if err := querier.LoadBankBalances(account, httpClient, chain); err != nil {
					bar.Describe(err.Error())
				}
				bar.Add(1)
				/*
					err = querier.LoadDistributionData(acct, httpClient)
					if err != nil {
						bar.Describe(err.Error())
					}
					bar.Add(1)

					err = querier.LoadStakingData(acct, httpClient)
					if err != nil {
						bar.Describe(err.Error())
					}
					bar.Add(1)*/
			}
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
			display.PrintAccountDetailsTable(&chains)
		}
	},
}

func init() {
	flagCsv = accountDetailsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	flagZbxAcctDetails = accountDetailsCmd.Flags().BoolP("zbx", "z", false, "push the result to a zabbix trapper item")
	accountsCmd.AddCommand(accountDetailsCmd)
}
