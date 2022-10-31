package config

import (
	"errors"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type AccountConfig struct {
	Name    string
	Address string
	Chain   string
}

type ChainConfig struct {
	ID  string
	LCD string
}

type Config struct {
	Accounts []AccountConfig
	Chains   []ChainConfig
}

func LoadConfig(configPath string) (model.Accounts, error) {

	if configPath != "" {
		p := filepath.Join(configPath)
		fmt.Println("p:", p)
		fmt.Println("Dir(p):", filepath.Dir(p))
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

	accounts := model.Accounts{}
	chains := model.Chains{}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if errors.Is(err, err.(viper.ConfigFileNotFoundError)) {
			if err != nil {
				if configPath != "" {
					viper.AddConfigPath(configPath)
					return accounts, errors.New("no configuration found at " + configPath)
				} else {
					return accounts, errors.New("no configuration found in default locations ($HOME/.stakooler) or (current directory)")
				}

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
