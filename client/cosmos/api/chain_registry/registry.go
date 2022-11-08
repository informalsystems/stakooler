package chain_registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AssetsList struct {
	Assets []struct {
		Description string `json:"description"`
		DenomUnits  []struct {
			Denom    string `json:"denom"`
			Exponent int    `json:"exponent"`
		} `json:"denom_units"`
		Base    string `json:"base"`
		Name    string `json:"name"`
		Display string `json:"display"`
		Symbol  string `json:"symbol"`
	} `json:"assets"`
}

func GetAssetsList(chain string) (AssetsList, error) {
	var assets AssetsList

	url := "https://chains.cosmos.directory/" + chain + "/assetlist"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return assets, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return assets, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return assets, err
	}
	err = json.Unmarshal(body, &assets)
	if err != nil {
		fmt.Println(err)
		return assets, err
	}
	return assets, nil
}
