package api

import (
	"encoding/json"
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"io/ioutil"
	"net/http"
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

func GetRewards(account *model.Account) (RewardsResponse, error) {
	var rewardsResponse RewardsResponse

	url := account.Chain.LCD + "/cosmos/distribution/v1beta1/delegators/" + account.Address + "/rewards"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return rewardsResponse, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return rewardsResponse, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return rewardsResponse, err
	}
	err = json.Unmarshal(body, &rewardsResponse)
	if err != nil {
		fmt.Println(err)
		return rewardsResponse, err
	}
	return rewardsResponse, nil
}
