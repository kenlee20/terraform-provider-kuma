package kuma

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (c *Client) GetNotifications() ([]Notification, error) {
	resp, err := c.doRequest("GET", "/notifications", nil)
	if err != nil {
		return nil, err
	}

	type _notifications struct {
		Notifications []Notification `json:"notifications"`
	}

	var notifications _notifications

	if err := json.Unmarshal(resp, &notifications); err != nil {
		return nil, err
	}

	return notifications.Notifications, nil
}

func (c *Client) GetNotification(id int) (*Notification, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/notifications/%s", strconv.Itoa(id)), nil)
	if err != nil {
		return nil, err
	}

	var notification Notification

	if err := json.Unmarshal(resp, &notification); err != nil {
		return nil, err
	}

	return &notification, nil
}

func (c *Client) GetDefaultNotifications() ([]int, error) {
	var defaultNotifications []int
	notifications, err := c.GetNotifications()
	if err != nil {
		return nil, err
	}

	for _, notification := range notifications {
		if notification.IsDefault {
			defaultNotifications = append(defaultNotifications, notification.ID)
		}
	}

	return defaultNotifications, nil
}
