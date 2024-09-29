package api

import (
	"encoding/json"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

type AcctResponse struct {
	Account struct {
		Type               string `json:"@type"`
		BaseVestingAccount struct {
			BaseAccount struct {
				Address       string `json:"address,omitempty"`
				PubKey        string `json:"public_key,omitempty"`
				AccountNumber string `json:"account_number,omitempty"`
				Sequence      string `json:"sequence,omitempty"`
			}
			OriginalVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"original_vesting"`
			DelegatedFree []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"delegated_free"`
			DelegatedVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"delegated_vesting"`
			EndTime string `json:"end_time"`
		} `json:"base_vesting_account"`
		StartTime      string `json:"start_time"`
		VestingPeriods []struct {
			Length string `json:"length"`
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"vesting_periods"`
	} `json:"account"`
}

func (a AcctResponse) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[OriginalVesting] = make(map[string]string)
	balances[DelegatedVesting] = make(map[string]string)

	for _, balance := range a.Account.BaseVestingAccount.OriginalVesting {
		balances[OriginalVesting][balance.Denom] = balance.Amount
	}

	for _, balance := range a.Account.BaseVestingAccount.DelegatedVesting {
		balances[DelegatedVesting][balance.Denom] = balance.Amount
	}
	return balances
}

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

func GetAccount(address string, endpoint string, client *http.Client) (response *AcctResponse, err error) {
	var body []byte

	url := endpoint + "/cosmos/auth/v1beta1/accounts/" + address
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &AcctResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return
	}
	return
}
