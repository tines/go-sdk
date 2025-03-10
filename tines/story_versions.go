package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type StoryVersion struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

type StoryVersionCreateRequest struct {
	Name    string `json:"name,omitempty"`
	DraftID int    `json:"draft_id,omitempty"`
}

// Create a new version of a story.
func (c *Client) CreateStoryVersion(ctx context.Context, storyID int, request *StoryVersionCreateRequest) (*StoryVersion, error) {
	resource := fmt.Sprintf("/api/v1/stories/%d/versions", storyID)
	storyVersion := StoryVersion{}

	req, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, http.MethodPost, resource, nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &storyVersion)
	if err != nil {
		return nil, err
	}

	return &storyVersion, nil
}
