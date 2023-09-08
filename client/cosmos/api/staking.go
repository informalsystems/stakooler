package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetDelegations(account *model.Account, client *http.Client) (response model.Delegations, err error) {
	var body []byte

	url := account.Chain.LCD + "/cosmos/staking/v1beta1/delegations/" + account.Address
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

func GetUnbondings(account *model.Account, client *http.Client) (response model.Unbondings, err error) {
	var body []byte

	url := account.Chain.LCD + "/cosmos/staking/v1beta1/delegators/" + account.Address + "/unbonding_delegations"
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

func GetStakingParams(chainEndpoint string, client *http.Client) (response model.Params, err error) {
	var body []byte

	url := chainEndpoint + "/cosmos/staking/v1beta1/params"
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

func GetChainValidators(validator *model.Validator) (model.Validators, error) {
	var validators model.Validators

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

func GetValidatorUnbondings(validator *model.Validator) (model.Unbondings, error) {
	var unbondings model.Unbondings
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

func GetValidatorDelegations(validator *model.Validator) (model.Delegations, error) {
	var delegations model.Delegations

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
