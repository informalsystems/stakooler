package cmd

import (
	"fmt"
	"os"

	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/api/chain_registry"
	"github.com/informalsystems/stakooler/client/cosmos/querier"
	"github.com/informalsystems/stakooler/client/display"
	"github.com/informalsystems/stakooler/config"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	flagCsvValidatorStats *bool
	flagZbxValidatorStats *bool
)

// represents the 'accounts details' command
var validatorStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Shows detailed information about a validator statistics",
	Long: `This command shows detailed information about a validator statistics. For example:

It shows the validator's voting power, voting power percentage, ranking, number of delegators per chain`,
	Run: func(cmd *cobra.Command, args []string) {
		barEnabled := !*flagCsvValidatorStats && !*flagZbxValidatorStats
		config, err := config.LoadConfig(flagConfigPath)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading configuration file")
			os.Exit(1)
		}

		if *flagZbxValidatorStats {
			if config.Zabbix.Port <= 0 || config.Zabbix.Host == "" {
				log.Fatal().Err(err).Msg("zabbix output requested. missing or incorrect zabbix configuration")
				os.Exit(1)
			}
		}

		var bar *progressbar.ProgressBar

		if barEnabled {
			// Progress bar
			// iterations are the api calls number times the number of accounts
			totalIterations := len(config.Validators.Entries) * 5 // two API calls
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
		for _, validator := range config.Validators.Entries {
			// Don't show this if csv option enabled
			if barEnabled {
				bar.Describe(fmt.Sprintf("Getting statistics for %s", validator.ValoperAddress))
			}

			// Load assets list
			// Get Assets list for the chain
			assets, err := chain_registry.GetAssetsList(validator.Chain.ID)
			if err != nil {
				log.Fatal().Err(err).Str("chain", validator.Chain.ID).Msg("cannot retrieve assets list")
				os.Exit(1)
			}

			for _, asset := range assets.Assets {
				denom, err := api.GetStakingParams(validator.Chain.LCD)
				if err != nil {
					log.Fatal().Err(err).Str("chain", validator.Chain.ID).Msg("cannot retrieve staking params")
					os.Exit(1)
				}
				if asset.Base == denom.ParamsResponse.BondDenom {
					validator.Chain.Denom = asset.Symbol
					for _, du := range asset.DenomUnits {
						if du.Denom == asset.Display {
							validator.Chain.Exponent = du.Exponent
						}
					}
				}
			}

			// Load validators
			err = querier.LoadValidatorStats(validator, bar)
			if err != nil {
				log.Error().Err(err).Str("validator_addr", validator.ValoperAddress).Msg("error loading validator stats")
			} else {
				// Don't show this if csv option enabled
				if barEnabled {
					bar.Describe(fmt.Sprintf("Got validator %s statistics", validator.ValoperAddress))
				}
			}
		}

		// Hide bar
		bar.Finish()

		// If csv flag specified use csv output
		if *flagCsvValidatorStats {
			// write csv file
			display.WriteValidatorCSV(&config.Validators)
		} else if !*flagZbxValidatorStats {
			// Print table information
			display.PrintValidatorStasTable(&config.Validators)
		} else {
			display.Send(config.Zabbix.Host, config.Zabbix.Port, &config.Validators)
		}
	},
}

func init() {
	flagCsvValidatorStats = validatorStatsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	flagZbxValidatorStats = validatorStatsCmd.Flags().BoolP("zbx", "z", false, "push the result to a zabbix trapper item")
	validatorCmd.AddCommand(validatorStatsCmd)
}
