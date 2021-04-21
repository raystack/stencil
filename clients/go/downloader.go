package stencil

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func downloader(uri string, opts HTTPOptions) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request. %w", err)
	}
	for key, val := range opts.Headers {
		req.Header.Add(key, val)
	}
	res, err := (&http.Client{Timeout: opts.Timeout}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed. %w", err)
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		return ioutil.ReadAll(res.Body)
	default:
		body, err := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("request failed. response body: %s, response_read_error: %w", body, err)
	}
}
