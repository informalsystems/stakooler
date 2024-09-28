package api

import (
	"encoding/json"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetPrefix(endpointURL string, client *http.Client) (response model.Bech32PrefixResponse, err error) {
	var body []byte

	url := endpointURL + "/cosmos/auth/v1beta1/bech32"
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

func GetAccount(address string, endpoint string, client *http.Client) (response *model.AcctResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/auth/v1beta1/accounts/" + address
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &model.AcctResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return
	}
	return
}
