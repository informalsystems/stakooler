package cosmos

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

func HttpGet(url string, client *http.Client) (body []byte, err error) {
	var req *http.Request
	var res *http.Response

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err = client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("failed http query: %d", req.Response.StatusCode))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, err = io.ReadAll(res.Body)

	return
}
