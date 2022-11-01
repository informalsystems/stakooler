package config

import (
	"errors"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"strings"
)

type AccountConfig struct {
	Name    string
	Address string
	Chain   string
}

type ValidatorsConfig struct {
	Name    string
	Address string
	Chains  []string
}

type ChainConfig struct {
	ID  string
	LCD string
}

type Configuration struct {
	Accounts   []AccountConfig
	Validators []ValidatorsConfig
	Chains     []ChainConfig
}

func LoadConfig(configPath string) (model.Config, error) {

	if configPath != "" {
		p := filepath.Join(configPath)
		filename := filepath.Base(p)
		ext := filepath.Ext(filename)
		configName := strings.TrimSuffix(filename, ext)
		path := filepath.Dir(p)
		viper.SetConfigName(configName)
		viper.SetConfigType(strings.Replace(ext, ".", "", 1))
		viper.AddConfigPath(path)
	} else {
		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("toml")
		viper.AddConfigPath("$HOME/.stakooler") // call multiple times to add many search paths
		viper.AddConfigPath(".")
	}

	config := model.Config{}
	accounts := model.Accounts{}
	chains := model.Chains{}
	validators := model.Validators{}

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		if reflect.TypeOf(err).Kind() == reflect.TypeOf(viper.ConfigFileNotFoundError{}).Kind() {
			if configPath != "" {
				return config, errors.New("no configuration found at " + configPath)
			} else {
				return config, errors.New("cannot find config.toml in default locations ($HOME/.stakooler) or (current directory)")
			}
		} else {
			return config, errors.New(fmt.Sprintf("%s", err))
		}
	} else {
		var configuration Configuration
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return config, errors.New(fmt.Sprintf("can not decode configuration: %s", err))
		}

		for chIdx := range configuration.Chains {
			chain := model.Chain{
				ID:  configuration.Chains[chIdx].ID,
				LCD: configuration.Chains[chIdx].LCD,
			}
			chains.Entries = append(chains.Entries, chain)
		}

		for accIdx := range configuration.Accounts {
			found := false
			for _, c := range chains.Entries {
				if strings.ToUpper(c.ID) == strings.ToUpper(configuration.Accounts[accIdx].Chain) {
					account := model.Account{
						Name:    configuration.Accounts[accIdx].Name,
						Address: configuration.Accounts[accIdx].Address,
						Chain:   c,
					}
					accounts.Entries = append(accounts.Entries, &account)
					found = true
				}
			}
			if !found {
				return config, errors.New(fmt.Sprintf("can not find chain id specified for account %s (%s) in the config", configuration.Accounts[accIdx].Name, configuration.Accounts[accIdx].Address))
			}
		}
		config.Accounts = accounts
		config.Validators = validators
		config.Chains = chains

		return config, nil
	}
}
