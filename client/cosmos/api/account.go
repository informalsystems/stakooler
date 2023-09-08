package api

import (
	"encoding/json"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetAuth(account *model.Account, client *http.Client) (response model.AuthResponse, err error) {
	var body []byte

	url := account.Chain.LCD + "/cosmos/auth/v1beta1/accounts/" + account.Address
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
