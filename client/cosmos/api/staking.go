package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

type Delegations struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Unbondings struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

func (d *Delegations) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[Delegation] = make(map[string]string)

	for _, balance := range d.DelegationResponses {
		balances[Delegation][balance.Balance.Denom] = balance.Balance.Amount
	}
	return balances
}

func (u *Unbondings) GetBalances() map[int]map[string]string {
	balances := make(map[int]map[string]string)
	balances[Unbonding] = make(map[string]string)

	for _, response := range u.UnbondingResponses {
		for _, entry := range response.Entries {
			balances[Unbonding]["denom"] = entry.Balance
		}
	}
	return balances
}

func GetDelegations(address string, endpoint string, client *http.Client) (response *Delegations, err error) {
	var body []byte

	url := endpoint + "/cosmos/staking/v1beta1/delegations/" + address
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &Delegations{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return
	}
	return
}

func GetUnbondings(address string, endpoint string, client *http.Client) (response *Unbondings, err error) {
	var body []byte

	url := endpoint + "/cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations"
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &Unbondings{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

func GetStakingParams(chainEndpoint string, client *http.Client) (response model.Params, err error) {
	var body []byte

	url := chainEndpoint + "/cosmos/staking/v1beta1/params"
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

func GetChainValidators(endpoint string) (model.Validators, error) {
	var validators model.Validators

	url := endpoint + "/cosmos/staking/v1beta1/validators?pagination.limit=1000&pagination.count_total=true&status=BOND_STATUS_BONDED"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return validators, err
	}
	res, err := client.Do(req)
	if err != nil {
		return validators, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return validators, err
	}
	err = json.Unmarshal(body, &validators)
	if err != nil {
		fmt.Println(err)
		return validators, err
	}

	// Add block height
	validators.BlockHeight = res.Header.Get("Grpc-Metadata-X-Cosmos-Block-Height")

	return validators, nil
}

func GetValidatorUnbondings(endpoint string, address string) (Unbondings, error) {
	var unbondings Unbondings
	url := endpoint + "/cosmos/staking/v1beta1/validators/" + address + "/unbonding_delegations"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return unbondings, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return unbondings, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return unbondings, err
	}
	err = json.Unmarshal(body, &unbondings)
	if err != nil {
		fmt.Println(err)
		return unbondings, err
	}
	return unbondings, nil
}

func GetValidatorDelegations(endpoint string, valoper string) (Delegations, error) {
	var delegations Delegations

	url := endpoint + "/cosmos/staking/v1beta1/validators/" + valoper + "/delegations?pagination.limit=15000&pagination.count_total=true"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return delegations, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return delegations, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return delegations, err
	}
	err = json.Unmarshal(body, &delegations)
	if err != nil {
		fmt.Println(err)
		return delegations, err
	}
	return delegations, nil
}
