package kuma

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type requestOption func(*requestOptions)

type requestOptions struct {
	ContentType string
}

func withContentType(contentType string) requestOption {
	return func(o *requestOptions) {
		o.ContentType = contentType
	}
}

func NewClient(host, username, password *string) (*Client, error) {
	clearHost := strings.TrimRight(*host, "/")

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    clearHost,
	}

	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) doRequest(method string, uri string, rb io.Reader, opts ...requestOption) ([]byte, error) {
	token := c.Token
	clearUri := strings.TrimLeft(uri, "/")

	// Default options
	options := requestOptions{
		ContentType: "application/json",
	}

	for _, opt := range opts {
		opt(&options)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.HostURL, clearUri), rb)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", options.ContentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		error_messgae := fmt.Sprintf("method %s API %s", method, clearUri)
		return nil, fmt.Errorf("status: %d, message: %s, body: %s", res.StatusCode, error_messgae, body)
	}

	return body, err
}
