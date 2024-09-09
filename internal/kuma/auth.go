package kuma

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	readerBody := strings.NewReader(fmt.Sprintf("username=%s&password=%s", c.Auth.Username, c.Auth.Password))

	body, err := c.doRequest("POST", "/login/access-token/", readerBody, withContentType("application/x-www-form-urlencoded"))
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
