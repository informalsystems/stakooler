package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetAuth(account *model.Account) (model.AuthResponse, error) {
	var authResponse model.AuthResponse

	url := account.Chain.LCD + "/cosmos/auth/v1beta1/accounts/" + account.Address
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return authResponse, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return authResponse, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return authResponse, err
	}
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		fmt.Println(err)
		return authResponse, err
	}
	return authResponse, nil
}
