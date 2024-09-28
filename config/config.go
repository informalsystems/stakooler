package config

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/informalsystems/stakooler/client/cosmos/query"
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

func ParseChainConfig(data *model.RawAccountData, httpClient *http.Client) query.Chains {
	var chains query.Chains
	for _, chain := range data.Chains {
		chainData := &query.Chain{
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

		params, err := api.GetStakingParams(chainData.RestEndpoint, httpClient)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("cannot get chain prefix, skipping chain: %s", chainData.Id))
		} else {
			chainData.BondDenom = params.ParamsResponse.BondDenom
		}

		for _, acct := range data.Accounts {
			acctIndex := slices.IndexFunc(chain.Accounts, func(c string) bool { return c == acct.Name })

			if acctIndex != -1 {
				decoded, err := hex.DecodeString(acct.Address)
				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("cannot decode hex encoded account %s", acct.Address))
				}

				encodedAddr, err := bech32.ConvertAndEncode(chainData.Bech32Prefix, decoded)
				if err != nil {
					log.Error().Err(err).Msg("cannot bech32 encode address, skipping account")
					continue
				}

				encodeValoper, err := bech32.ConvertAndEncode(chainData.Bech32Prefix+"valoper", decoded)
				if err != nil {
					log.Error().Err(err).Msg("cannot bech32 encode address, skipping account")
					continue
				}

				chainData.Accounts = append(chainData.Accounts, &model.Account{
					Name:    acct.Name,
					Address: encodedAddr,
					Valoper: encodeValoper,
				})
			}
		}
		chains.Entries = append(chains.Entries, chainData)
	}
	return chains
}
