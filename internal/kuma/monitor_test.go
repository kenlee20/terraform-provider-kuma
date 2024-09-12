package kuma

import (
	"testing"
)

func TestCreateMonitorTag(t *testing.T) {
	host := "http://localhost:8000/"
	usr := "admin"
	pwd := "admin"

	client, err := NewClient(&host, &usr, &pwd)
	if err != nil {
		t.Error(err)
	}

	tags := MonitorTag{
		Name:  "demo4",
		Value: "789",
	}

	err = client.CreateMonitorTag(24, tags)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteMonitorTag(t *testing.T) {
	host := "http://localhost:8000/"
	usr := "admin"
	pwd := "admin"

	client, err := NewClient(&host, &usr, &pwd)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteMonitorTag(52, MonitorTag{TagId: 3, Value: "123"})
	if err != nil {
		t.Error(err)
	}
}
