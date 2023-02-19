package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient struct {
	c        *http.Client
	apiToken string
	host     string
}

const (
	headerAuthorisation = "Authorization"
	headerContentType   = "Content-Type"

	contentTypeJSON = "application/json"
)

func (c *httpClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *httpClient) Put(path string, body any) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, c.host+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add(headerAuthorisation, fmt.Sprintf("Bearer %s", c.apiToken))
	req.Header.Add(headerContentType, contentTypeJSON)
	return c.c.Do(req)
}
