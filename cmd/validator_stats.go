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

// represents the 'accounts details' command
var validatorStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Shows detailed information about a validator statistics",
	Long: `This command shows detailed information about a validator statistics. For example:

It shows the validator's voting power, voting power percentage, ranking, number of delegators per chain`,
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
			totalIterations := len(config.Validators.Entries) * 4 // two API calls
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

			err := querier.LoadValidatorStats(validator, bar)
			if err != nil {
				bar.Describe(fmt.Sprintf("failed to retrieve statistics for %s: %s", validator.ValoperAddress, err))
			} else {
				// Don't show this if csv option enabled
				if barEnabled {
					bar.Describe(fmt.Sprintf("Got validator %s statistics", validator.ValoperAddress))
				}
			}
		}

		// Hide bar
		bar.Finish()

		// TODO: Implement outputs
		// If csv flag specified use csv output
		if *flagCsv {
			// write csv file
			display.WriteValidatorCSV(&config.Validators)
		} else {
			// Print table information
			display.PrintValidatorStasTable(&config.Validators)
		}
	},
}

func init() {
	flagCsv = validatorStatsCmd.Flags().BoolP("csv", "c", false, "output the result to a csv format")
	validatorCmd.AddCommand(validatorStatsCmd)
}
