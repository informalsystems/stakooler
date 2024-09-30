package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AssetList struct {
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

type ChainData struct {
	Bech32Prefix string `json:"bech32_prefix"`
}

// SearchForAsset search for the symbol for a particular denom in the assets list
func (a *AssetList) SearchForAsset(denom string) (string, int) {
	for i := range a.Assets {
		if a.Assets[i].Base == denom {
			for j := range a.Assets[i].DenomUnits {
				if strings.ToUpper(a.Assets[i].DenomUnits[j].Denom) == strings.ToUpper(a.Assets[i].Display) {
					// Some chains (*cough* Injective *cough*) have decided to differentiate between denom units
					// using letter casing...
					if a.Assets[i].DenomUnits[j].Exponent != 0 {
						return a.Assets[i].Symbol, a.Assets[i].DenomUnits[j].Exponent
					}
				}
			}
		}
	}
	return denom, 0
}

func (a *AssetList) GetAssetsList(chain string, client *http.Client) error {
	var body []byte

	url := "https://chains.cosmos.directory/" + chain + "/assetlist"
	body, err := HttpGet(url, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, a)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChainData) GetChainData(chain string, client *http.Client) error {
	var body []byte

	url := "https://chains.cosmos.directory/" + chain + "/chain"
	body, err := HttpGet(url, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}
	return nil
}
