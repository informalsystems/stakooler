/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stakooler",
	Short: "Stakooler is a helpful utility for Cosmos delegators",
	Long: `Stakooler is a utility that allows delegators to retrieve information about 
configured Cosmos chains accounts. 

The tool allows the information retrieval for all configured accounts with one command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	cobra.OnInitialize(initConfig)
	options := cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
	}
	rootCmd.CompletionOptions = options

	// Add flags
	addGlobalFlags(rootCmd)
}

func initConfig() {

}
