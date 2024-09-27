package api

import (
	"encoding/json"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetRewards(account *model.Account, client *http.Client) (response model.RewardsResponse, err error) {
	var body []byte

	url := account.Chain.RestEndpoint + "/cosmos/distribution/v1beta1/delegators/" + account.Address + "/rewards"
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

func GetCommissions(account *model.Account, validator string, client *http.Client) (response model.CommissionResponse, err error) {
	var body []byte

	url := account.Chain.RestEndpoint + "/cosmos/distribution/v1beta1/validators/" + validator + "/commission"
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
