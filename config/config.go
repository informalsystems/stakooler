package config

import (
	"errors"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/viper"
	"strings"
)

type AccountConfig struct {
	Name	string
	Address	string
	Chain	string
}

type ChainConfig struct {
	ID	string
	LCD	string
}

type Config struct {
	Accounts	[]AccountConfig
	Chains		[]ChainConfig
}

func LoadConfig() (model.Accounts, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.stakooler") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	accounts := model.Accounts{}
	chains := model.Chains{}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if errors.Is(err, err.(viper.ConfigFileNotFoundError)) {
			if err != nil {
				return accounts, errors.New("no configuration found")
			}
		} else {
			return accounts, errors.New(fmt.Sprintf("error loading configuration file: %s", err))
		}
	} else {
		var config Config
		err := viper.Unmarshal(&config)
		if err != nil {
			return accounts, errors.New(fmt.Sprintf("can not decode configuration: %s", err))
		}

		for chIdx := range config.Chains {
			chain := model.Chain{
				ID:  config.Chains[chIdx].ID,
				LCD: config.Chains[chIdx].LCD,
			}
			chains.Entries = append(chains.Entries, chain)
		}

		for accIdx := range config.Accounts {
			found := false
			for _, c := range chains.Entries {
				if strings.ToUpper(c.ID) == strings.ToUpper(config.Accounts[accIdx].Chain) {
					account := model.Account{
						Name:    config.Accounts[accIdx].Name,
						Address: config.Accounts[accIdx].Address,
						Chain:   c,
					}
					accounts.Entries = append(accounts.Entries, &account)
					found = true
				}
			}
			if !found {
				return accounts, errors.New(fmt.Sprintf("can not find chain id specified for account %s (%s) in the config", config.Accounts[accIdx].Name, config.Accounts[accIdx].Address))
			}
		}
		return accounts, nil
	}
	return accounts, nil
}
