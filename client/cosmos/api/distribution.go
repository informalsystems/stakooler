package api

import (
	"encoding/json"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetRewards(address string, endpoint string, client *http.Client) (response model.RewardsResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"
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

func GetCommissions(validator string, endpoint string, client *http.Client) (response model.CommissionResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/distribution/v1beta1/validators/" + validator + "/commission"
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
