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
		HTTPClient: &http.Client{Timeout: 300 * time.Second},
		HostURL:    clearHost,
		Retry:      10,
		Interval:   10 * time.Second,
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

func (c *Client) doRequest(method string, uri string, rb io.Reader, opts ...requestOption) ([]byte, *int, error) {
	var error_messgae string
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
		return nil, nil, err
	}

	req.Header.Add("Content-Type", options.ContentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		error_messgae = fmt.Sprintf("method %s API %s: %v", method, clearUri, err)
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, &res.StatusCode, readErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, &res.StatusCode, fmt.Errorf("status: %s, message: %s, body: %s", res.Status, error_messgae, body)
	}

	return body, &res.StatusCode, nil
}
