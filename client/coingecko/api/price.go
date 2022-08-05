package api

import (
	"fmt"
	"github.com/informalsystems/stakooler/utils"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PriceResponse struct {
	Token struct {
		Usd string `json:"usd"`
	} `json:""`
}

func GetTokenPrice(id string) (float64, error) {

	url := "https://api.coingecko.com/api/v3/simple/price?ids=" + id + "&vs_currencies=usd"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return utils.ZEROAMOUNT, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return utils.ZEROAMOUNT, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return utils.ZEROAMOUNT, err
	}
	value := gjson.Get(string(body), id+".usd")
	price, err := strconv.ParseFloat(value.Raw, 64)
	if err != nil {
		fmt.Println(err)
		return utils.ZEROAMOUNT, err
	}
	return price, nil
}
