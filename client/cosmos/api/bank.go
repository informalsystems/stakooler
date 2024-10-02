package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type BankResponse struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type DenomMetadataResponse struct {
	Metadata struct {
		Description string `json:"description"`
		DenomUnits  []struct {
			Denom    string   `json:"denom"`
			Exponent int      `json:"exponent"`
			Aliases  []string `json:"aliases"`
		} `json:"denom_units"`
		Base    string `json:"base"`
		Display string `json:"display"`
	} `json:"metadata"`
}

func (b *BankResponse) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[Bank] = map[string]string{}

	for _, balance := range b.Balances {
		balances[Bank][balance.Denom] = balance.Amount
	}
	return balances
}

func (b *BankResponse) QueryBankBalances(address string, endpoint string, client *http.Client) error {
	var body []byte

	url := endpoint + "/cosmos/bank/v1beta1/balances/" + address
	body, err := HttpGet(url, client)

	err = json.Unmarshal(body, b)
	if err != nil {
		return err
	}
	return nil
}

func (d *DenomMetadataResponse) QueryMetadataFromBank(denom string, endpoint string, client *http.Client) error {
	var body []byte
	var url string

	// This is here because of injective does not implement the original query
	if strings.HasPrefix(denom, "factory/inj") {
		url = endpoint + "/cosmos/bank/v1beta1/denoms_metadata_by_query_string?denom=" + denom
	} else {
		url = endpoint + "/cosmos/bank/v1beta1/denoms_metadata/" + denom
	}

	body, err := HttpGet(url, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}
	return err
}

func (d *DenomMetadataResponse) GetExponent() int {
	for _, unit := range d.Metadata.DenomUnits {
		if strings.ToUpper(unit.Denom) == strings.ToUpper(d.Metadata.Display) {
			return unit.Exponent
		}
	}
	return 0
}
