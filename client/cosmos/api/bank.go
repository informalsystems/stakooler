package api

import (
	"encoding/json"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"io/ioutil"
	"net/http"
)

const ApiBankBalances string = "/cosmos/bank/v1beta1/balances/"

type BalancesResponse struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

func GetBalances(account *model.Account) (BalancesResponse, error) {
	var balanceResponse BalancesResponse

	url := account.Chain.LCD + ApiBankBalances + account.Address
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return balanceResponse, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return balanceResponse, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return balanceResponse, err
	}
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		fmt.Println(err)
		return balanceResponse, err
	}
	return balanceResponse, nil
}
