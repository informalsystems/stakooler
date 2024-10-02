package cmd

import (
	"fmt"

	"github.com/informalsystems/stakooler/client/cosmos/api"
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
		rawAcctData, err := config.ReadAccountData("")
		if err != nil {
			log.Fatal().Err(err).Msg("error reading account data file")
		}

		var bar *progressbar.ProgressBar

		httpClient := api.NewHttpClient()
		chains := config.ParseAccountsConfig(rawAcctData, httpClient)

		if barEnabled {
			// Progress bar
			// iterations are the api calls number times the number of accounts
			totalIterations := len(chains)
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

		for _, chain := range chains {
			blockInfo := api.BlockResponse{}
			if err := blockInfo.GetLatestBlock(chain.RestEndpoint, httpClient); err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("failed to get latest block, skipping chain %s", chain.Id))
			}

			if barEnabled {
				bar.Describe(fmt.Sprintf("Getting chain %s details", chain.Id))
			}

			if err = chain.FetchAccountBalances(blockInfo, httpClient); err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("failed fetching accounts for %s", chain.Name))
			}
			bar.Add(1)
		}

		if err := bar.Finish(); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to finish bar"))
		}

		if *flagCsv {
			display.WriteAccountsCSV(chains)
		} else {
			display.PrintAccountDetailsTable(chains)
		}
	},
}

func init() {
	flagCsv = accountDetailsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	accountsCmd.AddCommand(accountDetailsCmd)
}
