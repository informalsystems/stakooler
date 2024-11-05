package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AssetPair struct {
	Rate         float64 `json:"rate"`
	AssetIdBase  string  `json:"asset_id_base"`
	AssetIdQuote string  `json:"asset_id_quote"`
}

func (c *AssetPair) GetCoinGekoQuote() error {

	var url string
	key := os.Getenv("COINAPI_KEY")

	if key == "" {
		return errors.New("COINAPI_KEY environment variable is not set")
	}

	if c.AssetIdQuote == "USD" {
		url = "https://api.coingecko.com/api/v3/simple/price?ids=" + c.AssetIdBase + "&vs_currencies=usd"
	} else {
		url = "https://api.coingecko.com/api/v3/simple/price?ids=" + c.AssetIdBase + "&vs_currencies=cad"
	}

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("coingeko request failed with status code: %s", res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *AssetPair) GetCoinApiQuote() error {

	var url string
	key := os.Getenv("COINAPI_KEY")

	if key == "" {
		return errors.New("COINAPI_KEY environment variable is not set")
	}

	if c.AssetIdQuote == "USD" {
		url = "https://rest.coinapi.io/v1/exchangerate/" + c.AssetIdBase + "/USD"
	} else {
		url = "https://rest.coinapi.io/v1/exchangerate/" + c.AssetIdBase + "/CAD"
	}

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "text/plain")
	req.Header.Add("X-CoinAPI-Key", key)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("coinapi request failed with status code: %s", res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}

	return nil
}
