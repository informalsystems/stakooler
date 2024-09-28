package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetBalances(address string, endpoint string, client *http.Client) (response model.BalancesResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/bank/v1beta1/balances/" + address
	body, err = cosmos.HttpGet(url, client)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

func GetDenomMetadataFromBank(denom string, endpoint string, client *http.Client) (response model.DenomMetadataResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/bank/v1beta1/denoms_metadata/" + denom
	body, err = cosmos.HttpGet(url, client)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

func GetExponent(metadata *model.DenomMetadataResponse) int {
	exponent := 0
	for _, d := range metadata.Metadata.DenomUnits {
		if strings.ToUpper(d.Denom) == strings.ToUpper(metadata.Metadata.Display) {
			return d.Exponent
		}
	}
	return exponent
}
