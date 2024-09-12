package kuma

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) GetTags() ([]Tag, error) {
	body, err := c.doRequest("GET", "/tags", nil)
	if err != nil {
		return nil, err
	}

	type _tags struct {
		Tags []Tag `json:"tags"`
	}

	var tags _tags

	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	return tags.Tags, nil
}

func (c *Client) GetTag(tagName string) (*Tag, error) {
	var t Tag

	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.Name == tagName {
			t = tag
			break
		}
	}
	return &t, nil
}

func (c *Client) CreateTag(tag Tag) (*Tag, error) {
	rb, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest("POST", "/tags", strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	newTag := Tag{}

	if err := json.Unmarshal(body, &newTag); err != nil {
		return nil, err
	}

	return &newTag, nil
}

func (c *Client) DeleteTag(tagId int64) error {
	uri := fmt.Sprintf("/tags/%s", strconv.FormatInt(tagId, 10))
	_, err := c.doRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateTag(tagId int64, tagInfo Tag) error {
	rb, err := json.Marshal(tagInfo)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PATCH", fmt.Sprintf("/tags/%s", strconv.FormatInt(tagId, 10)), strings.NewReader(string(rb)))

	return err
}
