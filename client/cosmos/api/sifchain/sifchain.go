package sifchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Token Registry information from Sifchain hosted in Github
type SifchainTokenRegistry struct {
	Tokens []struct {
		Decimals                 string   `json:"decimals"`
		Denom                    string   `json:"denom"`
		BaseDenom                string   `json:"base_denom"`
		Path                     string   `json:"path"`
		IbcChannelID             string   `json:"ibc_channel_id"`
		IbcCounterpartyChannelID string   `json:"ibc_counterparty_channel_id"`
		DisplayName              string   `json:"display_name"`
		DisplaySymbol            string   `json:"display_symbol"`
		Network                  string   `json:"network"`
		Address                  string   `json:"address"`
		ExternalSymbol           string   `json:"external_symbol"`
		TransferLimit            string   `json:"transfer_limit"`
		Permissions              []string `json:"permissions"`
		UnitDenom                string   `json:"unit_denom"`
		IbcCounterpartyDenom     string   `json:"ibc_counterparty_denom"`
		IbcCounterpartyChainID   string   `json:"ibc_counterparty_chain_id"`
	} `json:"entries"`
}

// Get Token List metadata hosted in Github and it is used by Sifchain
func GetTokenList() (SifchainTokenRegistry, error) {

	list := SifchainTokenRegistry{}
	url := "https://raw.githubusercontent.com/Sifchain/sifnode/develop/scripts/ibc/tokenregistration/sifchain-1/tokenregistry.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return list, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return list, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		fmt.Println(err)
		return list, err
	}
	return list, nil
}

// Searches for the symbol for a particular denom in the tokens list
// Returns the symbol and decimals precision
func (a *SifchainTokenRegistry) GetSymbolExponent(denom string) (string, int) {
	for i := range a.Tokens {
		if strings.ToUpper(a.Tokens[i].Denom) == strings.ToUpper(denom) {
			decimals, err := strconv.Atoi(a.Tokens[i].Decimals)
			if err != nil {
				return a.Tokens[i].Denom, 0
			} else {
				return a.Tokens[i].Denom, decimals
			}
		}
	}
	return denom, 0
}
