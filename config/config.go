package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	_ "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"slices"
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

func ReadAccountData(path string) (*model.RawAccountData, error) {
	var rawAcctData model.RawAccountData
	if path != "" {
		cleanPath := filepath.Join(path)
		fileName := filepath.Base(cleanPath)
		ext := filepath.Ext(fileName)
		trimmedName := strings.TrimSuffix(fileName, ext)
		filePath := filepath.Dir(cleanPath)
		viper.SetConfigName(trimmedName)
		viper.SetConfigType(strings.Replace(ext, ".", "", 1))
		viper.AddConfigPath(filePath)
	} else {
		viper.SetConfigName("accounts")
		viper.SetConfigType("json")
		viper.AddConfigPath("$HOME/.stakooler")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("error reading file")
		return nil, err
	} else {
		if err = viper.Unmarshal(&rawAcctData); err != nil {
			log.Error().Err(err).Msg("cannot unmarshall account data file")
			return nil, err
		}
	}
	return &rawAcctData, nil
}

func ParseChainConfig(data *model.RawAccountData, httpClient *http.Client) model.Chains {
	var chains model.Chains
	for _, chain := range data.Chains {
		chainData := &model.Chain{
			Name:         chain.Name,
			Id:           chain.Id,
			RestEndpoint: chain.Rest,
		}

		if prefix, err := api.GetPrefix(chainData.RestEndpoint, httpClient); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("cannot get chain prefix, skipping chain %s", chainData.Id))
			continue
		} else {
			chainData.Bech32Prefix = prefix.Bech32Prefix
		}

		for _, acct := range data.Accounts {
			acctIndex := slices.IndexFunc(chain.Accounts, func(c string) bool { return c == acct.Name })

			if acctIndex != -1 {
				decoded, err := hex.DecodeString(acct.Address)
				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("cannot decode hex encoded account %s", acct.Address))
				}
				if encodedAddr, err := bech32.ConvertAndEncode(chainData.Bech32Prefix, decoded); err != nil {
					log.Error().Err(err).Msg("cannot bech32 encode address, skipping account")
					continue
				} else {
					chainData.Accounts = append(chainData.Accounts, &model.Account{
						Name:    acct.Name,
						Address: encodedAddr,
					})
				}
			}
		}
		chains.Entries = append(chains.Entries, chainData)
	}
	return chains
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
	validators := model.ValidatorList{}

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
				Id:           configuration.Chains[chIdx].ID,
				RestEndpoint: configuration.Chains[chIdx].LCD,
				Denom:        configuration.Chains[chIdx].Denom,
				Exponent:     configuration.Chains[chIdx].Exponent,
			}
			chains.Entries = append(chains.Entries, &chain)
		}

		// Iterate through accounts in the configuration file
		for accIdx := range configuration.Accounts {
			found := false
			for _, c := range chains.Entries {
				if strings.ToUpper(c.Id) == strings.ToUpper(configuration.Accounts[accIdx].Chain) {
					account := model.Account{
						Name:    configuration.Accounts[accIdx].Name,
						Address: configuration.Accounts[accIdx].Address,
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
				if strings.ToUpper(c.Id) == strings.ToUpper(configuration.Validators[idx].Chain) {
					validator := model.Validator{
						ValoperAddress: configuration.Validators[idx].Valoper,
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
