package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetDelegations(address string, endpoint string, client *http.Client) (response *model.Delegations, err error) {
	var body []byte

	url := endpoint + "/cosmos/staking/v1beta1/delegations/" + address
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &model.Delegations{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return
	}
	return
}

func GetUnbondings(address string, endpoint string, client *http.Client) (response *model.Unbondings, err error) {
	var body []byte

	url := endpoint + "/cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations"
	body, err = HttpGet(url, client)
	if err != nil {
		return
	}

	response = &model.Unbondings{}
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

func GetValidatorUnbondings(endpoint string, address string) (model.Unbondings, error) {
	var unbondings model.Unbondings
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

func GetValidatorDelegations(endpoint string, valoper string) (model.Delegations, error) {
	var delegations model.Delegations

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
