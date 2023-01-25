package config

import (
	"errors"
	"fmt"
	_ "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type AccountConfig struct {
	Name    string
	Address string
	Chain   string
}

type ValidatorsConfig struct {
	Valoper string
	Chain   string
}

type ChainConfig struct {
	ID       string
	LCD      string
	Denom    string
	Exponent int
}

type Configuration struct {
	Accounts   []AccountConfig
	Validators []ValidatorsConfig
	Chains     []ChainConfig
	Zabbix     model.ZabbixConfig
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
		log.Error().Err(err).Msg("error reading configuration file")
		return config, err
	} else {
		var configuration Configuration
		err := viper.Unmarshal(&configuration)
		if err != nil {
			log.Error().Err(err).Msg("cannot unmarshall configuration file")
			return config, err
		}

		// Iterate through chains in the configuration file
		for chIdx := range configuration.Chains {
			chain := model.Chain{
				ID:       configuration.Chains[chIdx].ID,
				LCD:      configuration.Chains[chIdx].LCD,
				Denom:    configuration.Chains[chIdx].Denom,
				Exponent: configuration.Chains[chIdx].Exponent,
			}
			chains.Entries = append(chains.Entries, chain)
		}

		// Iterate through accounts in the configuration file
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
					break
				}
			}
			if !found {
				return config, errors.New(fmt.Sprintf("can not find chain id specified for account %s (%s) in the config", configuration.Accounts[accIdx].Name, configuration.Accounts[accIdx].Address))
			}
		}

		// Iterate through accounts in the configuration file
		for idx := range configuration.Validators {
			found := false
			for _, c := range chains.Entries {
				if strings.ToUpper(c.ID) == strings.ToUpper(configuration.Validators[idx].Chain) {
					validator := model.Validator{
						ValoperAddress: configuration.Validators[idx].Valoper,
						Chain:          c,
					}
					validators.Entries = append(validators.Entries, &validator)
					found = true
					break
				}
			}
			if !found {
				return config, errors.New(fmt.Sprintf("can not find chain id specified for account %s (%s) in the config", configuration.Accounts[idx].Name, configuration.Accounts[idx].Address))
			}
		}

		config.Zabbix.Host = configuration.Zabbix.Host
		config.Zabbix.Port = configuration.Zabbix.Port
		config.Zabbix.Server = configuration.Zabbix.Server
		config.Accounts = accounts
		config.Validators = validators
		config.Chains = chains

		return config, nil
	}
}
