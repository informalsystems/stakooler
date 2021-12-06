package api

import (
	"encoding/json"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"io/ioutil"
	"net/http"
)


type MintParams struct {
	Params struct {
		MintDenom           string `json:"mint_denom"`
	} `json:"params"`
}

func GetMintParams(account *model.Account) (MintParams, error) {
	var params MintParams

	url := account.Chain.LCD + "/cosmos/mint/v1beta1/params"
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

	body, err := ioutil.ReadAll(res.Body)
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
