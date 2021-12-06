package api

import (
	"encoding/json"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"io/ioutil"
	"net/http"
	"strings"
)

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

type DenomMetadataResponse struct {
	Metadata struct {
		Description string `json:"description"`
		DenomUnits  []struct {
			Denom    string   `json:"denom"`
			Exponent int      `json:"exponent"`
			Aliases  []string `json:"aliases"`
		} `json:"denom_units"`
		Base    string `json:"base"`
		Display string `json:"display"`
	} `json:"metadata"`
}

func GetBalances(account *model.Account) (BalancesResponse, error) {
	var balanceResponse BalancesResponse

	url := account.Chain.LCD + "/cosmos/bank/v1beta1/balances/" + account.Address
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

func GetDenomMetadata(account *model.Account, denom string) (DenomMetadataResponse, error) {
	var denomMetadata DenomMetadataResponse

	url := account.Chain.LCD + "/cosmos/bank/v1beta1/denoms_metadata/" + denom
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return denomMetadata, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return denomMetadata, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return denomMetadata, err
	}
	err = json.Unmarshal(body, &denomMetadata)
	if err != nil {
		fmt.Println(err)
		return denomMetadata, err
	}
	return denomMetadata, nil
}

func (metadata *DenomMetadataResponse) GetExponent() int {
	exponent := 0
	for _, d := range metadata.Metadata.DenomUnits {
		if strings.ToUpper(d.Denom) == strings.ToUpper(metadata.Metadata.Display) {
			return d.Exponent
		}
	}
	return exponent
}