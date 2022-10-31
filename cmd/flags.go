package cmd

import "github.com/spf13/cobra"

var (
	flagConfigPath string
)

// addGlobalFlags defines flags to be used regardless of the command used
func addGlobalFlags(cmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, "config", "f", "", "configuration file")
}
