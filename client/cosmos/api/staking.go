package api

import (
	"encoding/json"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"io"
	"net/http"
	"time"
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

type Params struct {
	ParamsResponse struct {
		UnbondingTime     string `json:"unbonding_time"`
		MaxValidators     int    `json:"max_validators"`
		MaxEntries        int    `json:"max_entries"`
		HistoricalEntries int    `json:"historical_entries"`
		BondDenom         string `json:"bond_denom"`
		MinCommissionRate string `json:"min_commission_rate"`
	} `json:"params"`
}

type ValidatorSet struct {
	BlockHeight string `json:"block_height"`
	Validators  []struct {
		Address string `json:"address"`
		PubKey  struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"pub_key"`
		VotingPower      string `json:"voting_power"`
		ProposerPriority string `json:"proposer_priority"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Validators struct {
	BlockHeight        string `json:"block_height,omitempty"`
	ValidatorsResponse []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

func GetDelegations(account *model.Account) (Delegations, error) {
	var delegations Delegations

	url := account.Chain.LCD + "/cosmos/staking/v1beta1/delegations/" + account.Address
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

func GetUnbondings(account *model.Account) (Unbondings, error) {
	var unbondings Unbondings

	url := account.Chain.LCD + "/cosmos/staking/v1beta1/delegators/" + account.Address + "/unbonding_delegations"
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

func GetStakingParams(chainEndpoint string) (Params, error) {
	var params Params

	url := chainEndpoint + "/cosmos/staking/v1beta1/params"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return params, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return params, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return params, err
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		fmt.Println(err)
		return params, err
	}
	return params, nil
}

//func GetValidatorSet(validator *model.Validator) (ValidatorSet, error) {
//	var valset ValidatorSet
//
//	url := validator.Chain.LCD + "/cosmos/base/tendermint/v1beta1/validatorsets/latest?pagination.limit=300&pagination.count_total=true"
//	method := "GET"
//
//	client := &http.Client{}
//	req, err := http.NewRequest(method, url, nil)
//
//	if err != nil {
//		return valset, err
//	}
//	res, err := client.Do(req)
//	if err != nil {
//		return valset, err
//	}
//	defer res.Body.Close()
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		return valset, err
//	}
//	err = json.Unmarshal(body, &valset)
//	if err != nil {
//		fmt.Println(err)
//		return valset, err
//	}
//	return valset, nil
//}

func GetValidators(validator *model.Validator) (Validators, error) {
	var validators Validators

	url := validator.Chain.LCD + "/cosmos/staking/v1beta1/validators?pagination.limit=1000&pagination.count_total=true&status=BOND_STATUS_BONDED"
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

func GetValidatorUnbondings(validator *model.Validator) (Unbondings, error) {
	var unbondings Unbondings

	url := validator.Chain.LCD + "/cosmos/staking/v1beta1/validators/" + validator.ValoperAddress + "/unbonding_delegations"
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

func GetValidatorDelegations(validator *model.Validator) (Delegations, error) {
	var delegations Delegations

	url := validator.Chain.LCD + "/cosmos/staking/v1beta1/validators/" + validator.ValoperAddress + "/delegations?pagination.limit=15000&pagination.count_total=true"
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
