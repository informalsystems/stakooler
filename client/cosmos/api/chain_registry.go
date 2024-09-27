package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"

	"github.com/informalsystems/stakooler/client/cosmos/model"
)

func GetAssetsList(chain string, client *http.Client) (*model.AssetList, error) {
	var assets model.AssetList

	url := "https://chains.cosmos.directory/" + chain + "/assetlist"
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("could not create asset list request")
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error when making asset list request")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Error().Err(err).Msg(fmt.Sprintf("error when making asset list request, status: %d", res.StatusCode))
		return nil, err
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
		return nil, err
	}

	err = json.Unmarshal(body, &assets)
	if err != nil {
		log.Error().Err(err).Msg("error when unmarshalling asset list body")
		return nil, err
	}
	return &assets, nil
}
