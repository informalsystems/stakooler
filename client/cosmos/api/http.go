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

	// returning the body even if the status code is not 200
	// allows the caller to decide if it should stop
	if res.StatusCode != http.StatusOK {
		return body, errors.New(fmt.Sprintf("status code: %d", res.StatusCode))
	}
	return body, nil
}
