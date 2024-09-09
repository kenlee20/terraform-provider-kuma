package kuma

import (
	"encoding/json"
	"strconv"
	"strings"
)

func (c *Client) GetMonitors() ([]Monitor, error) {
	body, err := c.doRequest("GET", "/monitors", nil)
	if err != nil {
		return nil, err
	}

	type _monitors struct {
		Monitors []Monitor `json:"monitors"`
	}

	var monitors _monitors

	if err := json.Unmarshal(body, &monitors); err != nil {
		return nil, err
	}

	return monitors.Monitors, nil
}

func (c *Client) GetMonitor(id int) (*Monitor, error) {
	body, err := c.doRequest("GET", "/monitors/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	type _monitor struct {
		Monitor Monitor `json:"monitor"`
	}

	var monitor _monitor

	if err := json.Unmarshal(body, &monitor); err != nil {
		return nil, err
	}

	return &monitor.Monitor, nil
}

func (c *Client) CreateMonitor(monitor Monitor) (*int, error) {
	rb, err := json.Marshal(monitor)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest("POST", "/monitors", strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	type _resp struct {
		Msg       string `json:"msg"`
		MonitorID int    `json:"monitorId"`
	}

	resp := _resp{}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.MonitorID, nil
}

func (c *Client) DeleteMonitor(id int) error {
	_, err := c.doRequest("DELETE", "/monitors/"+strconv.Itoa(id), nil)
	return err
}

func (c *Client) UpdateMonitor(monitorID int, monitor Monitor) error {
	rb, err := json.Marshal(monitor)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PATCH", "/monitors/"+strconv.Itoa(monitorID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	return nil
}
