package model

import "strings"

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

// GetAssetDetails search for the symbol for a particular denom in the assets list
func (a *AssetList) GetAssetDetails(denom string) (string, int) {
	for i := range a.Assets {
		if a.Assets[i].Base == denom {
			for j := range a.Assets[i].DenomUnits {
				if strings.ToUpper(a.Assets[i].DenomUnits[j].Denom) == strings.ToUpper(a.Assets[i].Display) {
					return a.Assets[i].Symbol, a.Assets[i].DenomUnits[j].Exponent
				}
			}
		}
	}
	return denom, 0
}
