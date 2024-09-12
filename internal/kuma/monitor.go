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

func (c *Client) GetMonitor(id int64) (*Monitor, error) {
	body, err := c.doRequest("GET", "/monitors/"+strconv.FormatInt(id, 10), nil)
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

func (c *Client) CreateMonitor(monitor Monitor) (*int64, error) {
	defaultNotification, err := c.GetDefaultNotifications()
	if err != nil {
		return nil, err
	}

	monitor.NotificationIDList = append(monitor.NotificationIDList, defaultNotification...)

	// Marshal the monitor
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
		MonitorID int64  `json:"monitorId"`
	}

	resp := _resp{}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.MonitorID, nil
}

func (c *Client) DeleteMonitor(id int64) error {
	_, err := c.doRequest("DELETE", "/monitors/"+strconv.FormatInt(id, 10), nil)
	return err
}

func (c *Client) UpdateMonitor(monitorID int64, monitor Monitor) error {
	rb, err := json.Marshal(monitor)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PATCH", "/monitors/"+strconv.FormatInt(monitorID, 10), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateMonitorTag(monitorID int64, tagSet MonitorTag) (err error) {
	tagSetup := make(map[string]any)

	tag, err := c.GetTag(tagSet.Name)
	if err != nil {
		return err
	}

	tagSetup["tag_id"] = tag.ID
	tagSetup["value"] = tagSet.Value

	rb, err := json.Marshal(tagSetup)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/monitors/"+strconv.FormatInt(monitorID, 10)+"/tag", strings.NewReader(string(rb)))

	return err
}

func (c *Client) DeleteMonitorTag(monitorID int64, tagSet MonitorTag) (err error) {
	tagSetup := make(map[string]any)

	tagSetup["tag_id"] = tagSet.TagId
	tagSetup["value"] = tagSet.Value

	tag, err := json.Marshal(tagSetup)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/monitors/"+strconv.FormatInt(monitorID, 10)+"/tag/", strings.NewReader(string(tag)))

	return err
}
