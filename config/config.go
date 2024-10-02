package config

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/informalsystems/stakooler/client/cosmos/api"
	"github.com/informalsystems/stakooler/client/cosmos/model"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

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

func ParseAccountsConfig(data *model.RawAccountData, httpClient *http.Client) []*model.Chain {
	var chains []*model.Chain
	for _, chain := range data.Chains {
		chainData := &model.Chain{
			Name:         chain.Name,
			Id:           chain.Id,
			RestEndpoint: chain.Rest,
			AssetList:    &api.AssetList{},
		}

		if err := chainData.AssetList.QueryAssetList(chain.Name, httpClient); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("query asset list: %s", chain.Id))
		}

		prefixResponse := api.Bech32PrefixResponse{}
		if err := prefixResponse.GetPrefix(chainData.RestEndpoint, httpClient); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("query chain prefix, trying asset list %s", chainData.Id))
			chainDataRegistry := api.ChainData{}
			if err = chainDataRegistry.QueryChainData(chain.Name, httpClient); err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("query chain data, skipping chain: %s", chainData.Id))
				continue
			} else {
				chainData.Bech32Prefix = chainDataRegistry.Bech32Prefix
			}
		} else if chain.Name == "panacea" {
			// mediblock incorrectly shows cosmos as the prefix, so here we are, setting it manually
			chainData.Bech32Prefix = "panacea"
		} else {
			chainData.Bech32Prefix = prefixResponse.Bech32Prefix
		}

		params := &api.StakingParamsResponse{}
		if err := params.QueryParams(chainData.RestEndpoint, httpClient); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("query staking paramas, skipping chain: %s", chainData.Id))
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
					Tokens:  make(map[string]*model.Token),
				})
			}
		}
		chains = append(chains, chainData)
	}
	return chains
}
