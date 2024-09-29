package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
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

// SearchForAsset search for the symbol for a particular denom in the assets list
func (a *AssetList) SearchForAsset(denom string) (string, int) {
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

func (a *AssetList) GetAssetsList(chain string, client *http.Client) error {

	url := "https://chains.cosmos.directory/" + chain + "/assetlist"
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("could not create asset list request")
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error when making asset list request")
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Error().Err(err).Msg(fmt.Sprintf("error when making asset list request, status: %d", res.StatusCode))
		return err
	}

	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			log.Error().Err(err2).Msg("error when closing asset list response body")
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("error when reading asset list body")
		return err
	}

	err = json.Unmarshal(body, a)
	if err != nil {
		log.Error().Err(err).Msg("error when unmarshalling asset list body")
		return err
	}
	return nil
}
