package osmosis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AssetsList struct {
	ChainID string `json:"chain_id"`
	Assets  []struct {
		Description string `json:"description,omitempty"`
		DenomUnits  []struct {
			Denom    string        `json:"denom"`
			Exponent int           `json:"exponent"`
			Aliases  []interface{} `json:"aliases"`
		} `json:"denom_units"`
		Base     string `json:"base"`
		Name     string `json:"name"`
		Display  string `json:"display"`
		Symbol   string `json:"symbol"`
		LogoURIs struct {
			Png string `json:"png"`
			Svg string `json:"svg"`
		} `json:"logo_URIs,omitempty"`
		CoingeckoID string `json:"coingecko_id,omitempty"`
		Ibc         struct {
			SourceChannel string `json:"source_channel"`
			DstChannel    string `json:"dst_channel"`
			SourceDenom   string `json:"source_denom"`
		} `json:"ibc,omitempty"`
	} `json:"assets"`
}

// Get Assets List metadata hosted in Github and it is used by Osmosis
func GetAssetsList() (AssetsList, error) {

	assetsList := AssetsList{}
	url := "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return assetsList, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return assetsList, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return assetsList, err
	}
	err = json.Unmarshal(body, &assetsList)
	if err != nil {
		fmt.Println(err)
		return assetsList, err
	}
	return assetsList, nil
}

// Searches for the symbol for a particular denom in the assets list
// Returns the symbol and a boolean (true if found, or false if not)
func (a *AssetsList) GetSymbol(denom string) (string, bool) {
	for i := range a.Assets {
		if strings.ToUpper(a.Assets[i].Base) == strings.ToUpper(denom) {
			return a.Assets[i].Symbol, true
		}
	}
	return denom, false
}

// Search for the exponent value based on the symbol. Look for value in
// the denom units of the assets list
func (a *AssetsList) GetExponent(symbol string) int {
	for i := range a.Assets {
		for j := range a.Assets[i].DenomUnits {
			if strings.ToUpper(a.Assets[i].DenomUnits[j].Denom) == strings.ToUpper(symbol) {
				return a.Assets[i].DenomUnits[j].Exponent
			}
		}
	}
	return 0
}
