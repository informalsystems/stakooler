package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewHttpClient() (client *http.Client) {
	client = &http.Client{Timeout: 10 * time.Second}
	return
}

func HttpGet(url string, client *http.Client) ([]byte, error) {
	var req *http.Request
	var res *http.Response
	var body []byte

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed http query: %d", res.StatusCode))
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
