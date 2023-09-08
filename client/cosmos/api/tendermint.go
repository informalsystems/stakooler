package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos"
	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetLatestBlock(chain model.Chain, client *http.Client) (response model.BlockResponse, err error) {
	var body []byte

	url := chain.LCD + "/cosmos/base/tendermint/v1beta1/blocks/latest"
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

func GetBlock(height string, chain model.Chain) (model.BlockResponse, error) {
	var response model.BlockResponse

	url := chain.LCD + "/cosmos/base/tendermint/v1beta1/blocks/" + height
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return response, err
	}
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
