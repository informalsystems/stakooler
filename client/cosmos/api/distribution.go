package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RewardsResponse struct {
	Rewards []struct {
		ValidatorAddress string `json:"validator_address"`
		Reward           []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"reward"`
	} `json:"rewards"`
	Total []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
}

type CommissionResponse struct {
	Commissions struct {
		Commission []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"commission"`
	} `json:"commission"`
}

func (r *RewardsResponse) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[Rewards] = map[string]string{}

	for _, rewards := range r.Rewards {
		for _, reward := range rewards.Reward {
			balances[Rewards][reward.Denom] = reward.Amount
		}
	}
	return balances
}

func (c *CommissionResponse) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[Commission] = map[string]string{}

	for _, commission := range c.Commissions.Commission {
		balances[Commission][commission.Denom] = commission.Amount
	}
	return balances
}

func (r *RewardsResponse) QueryRewards(address string, endpoint string, client *http.Client) error {
	var body []byte

	url := endpoint + "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"
	body, err := HttpGet(url, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	return err
}

func (c *CommissionResponse) QueryCommission(validator string, endpoint string, client *http.Client) error {
	var body []byte

	url := endpoint + "/cosmos/distribution/v1beta1/validators/" + validator + "/commission"
	body, err := HttpGet(url, client)
	if err != nil {
		if strings.Contains(string(body), "validator does not exist") {
			return nil
		} else {
			return err
		}
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}
	return err
}
