package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type StoryImportMode string

const (
	StoryModeNew     StoryImportMode = "new"
	StoryModeReplace StoryImportMode = "versionReplace"
)

type StoryImportRequest struct {
	NewName  string                 `json:"new_name"`
	Data     map[string]interface{} `json:"data"`
	TeamID   int                    `json:"team_id"`
	FolderID int                    `json:"folder_id,omitempty"`
	Mode     StoryImportMode        `json:"mode"`
}

type Story struct {
	ID                   int      `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	UserID               int      `json:"user_id,omitempty"`
	Description          string   `json:"description,omitempty"`
	KeepEventsFor        int      `json:"keep_events_for,omitempty"`
	Disabled             bool     `json:"disabled,omitempty"`
	Priority             bool     `json:"priority,omitempty"`
	STSEnabled           bool     `json:"send_to_story_enabled,omitempty"`
	STSAccessSource      string   `json:"send_to_story_access_source,omitempty"`
	STSAccess            string   `json:"send_to_story_access,omitempty"`
	STSSkillConfirmation bool     `json:"send_to_story_skill_use_requires_confirmation,omitempty"`
	SharedTeamSlugs      []string `json:"shared_team_slugs,omitempty"`
	EntryAgentID         int      `json:"entry_agent_id,omitempty"`
	ExitAgents           []int    `json:"exit_agents,omitempty"`
	TeamID               int      `json:"team_id,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	Guid                 string   `json:"guid,omitempty"`
	Slug                 string   `json:"slug,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	EditedAt             string   `json:"edited_at,omitempty"`
	Mode                 string   `json:"mode,omitempty"`
	FolderID             int      `json:"folder_id,omitempty"`
	Published            bool     `json:"published,omitempty"`
	ChangeControlEnabled bool     `json:"change_control_enabled,omitempty"`
	Locked               bool     `json:"locked,omitempty"`
	Owners               []int    `json:"owners,omitempty"`
}

type StoryList struct {
	Stories []Story       `json:"stories,omitempty"`
	Meta    paginate.Meta `json:"meta,omitempty"`
}

// Create a new story with an empty storyboard. For managing storyboard contents via
// API, using the ImportStory() function is the recommended approach.
func (c *Client) CreateStory(ctx context.Context, s *Story) (*Story, error) {
	newStory := Story{}

	req, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, "POST", "/api/v1/stories", nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &newStory)
	if err != nil {
		return nil, err
	}

	return &newStory, nil
}

// Get current state for a story.
func (c *Client) GetStory(ctx context.Context, id int) (story *Story, e error) {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	body, err := c.doRequest(ctx, "GET", resource, nil, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &story)
	if err != nil {
		return nil, err
	}

	return story, nil
}

// Update a story.
func (c *Client) UpdateStory(ctx context.Context, id int, values *Story) (*Story, error) {
	updatedStory := Story{}
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	req, err := json.Marshal(&values)
	if err != nil {
		return &updatedStory, err
	}

	body, err := c.doRequest(ctx, "PUT", resource, nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &updatedStory)
	if err != nil {
		return nil, err
	}

	return &updatedStory, nil
}

// Yields an iterator that returns individual Stories, optionally filtered by Team ID, story
// status, or Folder ID. If no other filters are specified, ListStories() will recurse
// through all pages of results until no more are available. If `filters.WithMaxResults()` is
// set, this function will yield either the actual set of results or the specified maximum
// number of results, whichever is less.
//
// Example Usage:
//
//	for s, err := range ListStories(ctx, WithMaxResults(10)) {
//		if err != nil {
//			...
//		}
//		fmt.Println(s.Name)
//	}
func (c *Client) ListStories(ctx context.Context, f ListFilter) iter.Seq2[Story, error] {
	var storyList, resultList StoryList
	resource := "/api/v1/stories"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	c.logger.Debug(fmt.Sprintf("max results requested: %d", page.TotalRequested))

	return func(yield func(Story, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(Story{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(Story{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.Stories {
				storyList.Stories = append(storyList.Stories, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = StoryList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range storyList.Stories {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// Delete a story.
func (c *Client) DeleteStory(ctx context.Context, id int) error {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)

	return err
}

// Delete a batch of multiple stories via a list of Story IDs.
func (c *Client) BatchDeleteStories(ctx context.Context, ids []int) error {
	resource := "/api/v1/stories/batch"
	payload := make(map[string][]int)
	payload["ids"] = ids

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, http.MethodDelete, resource, nil, data)

	return err
}

// Export the storyboard contents and metadata to JSON. By default, URLs (for example,
// webhook URLs) will not be randomized in the export, so anyone with access to the
// exported JSON will be able to identify and call those webhooks. If you are exporting
// a story for sharing or public consumption, we strongly recommend randomizing the URLs.
func (c *Client) ExportStory(ctx context.Context, id int, randomizeUrls bool) (map[string]interface{}, error) {
	resource := fmt.Sprintf("/api/v1/stories/%d/export", id)

	params := make(map[string]any)
	if randomizeUrls {
		params["randomize_urls"] = randomizeUrls
	} else {
		params = nil
	}

	export := make(map[string]interface{})

	res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &export)
	if err != nil {
		return nil, err
	}

	return export, nil
}

// Import a new story, or override an existing one.
func (c *Client) ImportStory(ctx context.Context, story *StoryImportRequest) (*Story, error) {
	newStory := Story{}

	req, err := json.Marshal(&story)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, "POST", "/api/v1/stories/import", nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &newStory)
	if err != nil {
		return nil, err
	}

	return &newStory, nil
}
