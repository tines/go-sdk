package tines_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

const (
	// Sanitized API response value as of 2025-02-03.
	testStoryResp = `
{
    "name": "New Story",
    "user_id": 1,
    "description": null,
    "keep_events_for": 604800,
    "disabled": false,
    "priority": false,
    "send_to_story_enabled": false,
    "send_to_story_access_source": "OFF",
    "send_to_story_access": null,
    "shared_team_slugs": [],
    "entry_agent_id": null,
    "exit_agents": [],
    "send_to_story_skill_use_requires_confirmation": true,
    "team_id": 1,
    "tags": [],
    "guid": "a72744c526e7d3e5b608f130a583c98b",
    "slug": "new_story",
    "created_at": "2025-02-03T00:00:00Z",
    "updated_at": "2025-02-03T00:00:00Z",
    "edited_at": "2025-02-03T00:00:00Z",
    "mode": "LIVE",
    "id": 1,
    "folder_id": null,
    "published": true,
    "change_control_enabled": false,
    "locked": false,
    "owners": []
}`
	// Sanitized API response value as of 2025-02-03.
	testUpdateStoryResp = `
{
    "name": "New Story",
    "user_id": 1,
    "description": "Description",
    "keep_events_for": 604800,
    "disabled": false,
    "priority": false,
    "send_to_story_enabled": false,
    "send_to_story_access_source": "OFF",
    "send_to_story_access": null,
    "shared_team_slugs": [],
    "entry_agent_id": null,
    "exit_agents": [],
    "send_to_story_skill_use_requires_confirmation": true,
    "team_id": 1,
    "tags": [],
    "guid": "a72744c526e7d3e5b608f130a583c98b",
    "slug": "new_story",
    "created_at": "2025-02-03T00:00:00Z",
    "updated_at": "2025-02-03T00:00:00Z",
    "edited_at": "2025-02-03T00:00:00Z",
    "mode": "LIVE",
    "id": 1,
    "folder_id": null,
    "published": true,
    "change_control_enabled": false,
    "locked": false,
    "owners": []
}`
	// Sanitized API response value as of 2025-02-03.
	testListStoriesResp = `
{
    "stories": [
        {
            "name": "Test Story",
            "user_id": 1,
            "description": null,
            "keep_events_for": 604800,
            "disabled": false,
            "priority": false,
            "send_to_story_enabled": false,
            "send_to_story_access_source": "OFF",
            "send_to_story_access": null,
            "shared_team_slugs": [],
            "entry_agent_id": null,
            "exit_agents": [],
            "send_to_story_skill_use_requires_confirmation": true,
            "team_id": 1,
            "tags": [],
            "guid": "a72744c526e7d3e5b608f130a583c98b",
            "slug": "new_story",
			"created_at": "2025-02-03T00:00:00Z",
			"updated_at": "2025-02-03T00:00:00Z",
			"edited_at": "2025-02-03T00:00:00Z",
            "mode": "LIVE",
            "id": 1,
            "folder_id": null,
            "published": true,
            "change_control_enabled": false,
            "locked": false,
            "owners": []
        }
    ],
    "meta": {
        "current_page": "https://example.tines.com/api/v1/stories?per_page=20&page=1",
        "previous_page": null,
        "next_page": null,
        "next_page_number": null,
        "per_page": 20,
        "pages": 1,
        "count": 1
    }
}`
	// Sanitized API response value as of 2025-02-03.
	testExportStoryResp = `
{
    "schema_version": 23,
    "standard_lib_version": 63,
    "action_runtime_version": 15,
    "name": "Test Story",
    "description": null,
    "guid": "3ef721e341e953727b057d4bd7bd65eb",
    "slug": "test_story",
    "agents": [
      {
        "type": "Agents::WebhookAgent",
        "name": "Webhook Action",
        "disabled": false,
        "description": null,
        "guid": "7977a7af73df9e18234ae8acb814d4fa",
        "origin_story_identifier": "cloud:8eeb4fdfeb8e2a573e55856936a8cf7f:3ef721e341e953727b057d4bd7bd65eb",
        "options": {
          "path": "2a5b9cd5b43d06329fd72f70c7cbeede",
          "secret": "cf881382af21ef97840c36aa9391f6cc",
          "verbs": "get,post"
        },
        "reporting": {
          "time_saved_value": 0,
          "time_saved_unit": "minutes"
        },
        "monitoring": {
          "monitor_all_events": false,
          "monitor_failures": false,
          "monitor_no_events_emitted": null
        },
        "template": {
          "created_from_template_guid": null,
          "created_from_template_version": null,
          "template_tags": []
        },
        "width": null
      },
      {
        "type": "Agents::EventTransformationAgent",
        "name": "Event Transform Action",
        "disabled": false,
        "description": null,
        "guid": "6bba07417c1732ae4b2e7b642dcc9cea",
        "origin_story_identifier": "cloud:8eeb4fdfeb8e2a573e55856936a8cf7f:3ef721e341e953727b057d4bd7bd65eb",
        "options": {
          "mode": "message_only",
          "loop": false,
          "payload": {
            "message": "This is an automatically generated message from Tines"
          }
        },
        "reporting": {
          "time_saved_value": 0,
          "time_saved_unit": "minutes"
        },
        "monitoring": {
          "monitor_all_events": false,
          "monitor_failures": false,
          "monitor_no_events_emitted": null
        },
        "template": {
          "created_from_template_guid": null,
          "created_from_template_version": null,
          "template_tags": []
        },
        "width": null,
        "schedule": null
      }
    ],
    "diagram_notes": [],
    "links": [
      {
        "source": 0,
        "receiver": 1
      }
    ],
    "diagram_layout": "{\"7977a7af73df9e18234ae8acb814d4fa\":[360,135],\"6bba07417c1732ae4b2e7b642dcc9cea\":[360,240]}",
    "send_to_story_enabled": false,
    "entry_agent_guid": null,
    "exit_agent_guids": [],
    "exit_agent_guid": null,
    "api_entry_action_guids": [],
    "api_exit_action_guids": [],
    "keep_events_for": 86400,
    "reporting_status": true,
    "send_to_story_access": null,
    "story_library_metadata": {},
    "parent_only_send_to_story": false,
    "monitor_failures": false,
    "send_to_stories": [],
    "synchronous_webhooks_enabled": false,
    "send_to_story_access_source": 0,
    "send_to_story_skill_use_requires_confirmation": true,
    "pages": [],
    "tags": [],
    "time_saved_unit": "minutes",
    "time_saved_value": 0,
    "origin_story_identifier": "cloud:8eeb4fdfeb8e2a573e55856936a8cf7f:3ef721e341e953727b057d4bd7bd65eb",
    "integration_product": null,
    "integration_vendor": null,
    "llm_product_instructions": "",
    "exported_at": "2025-01-01T00:00:00Z",
    "icon": ":house_buildings:",
    "integrations": []
  }`
	// Sanitized API response value as of 2025-02-03.
	testImportStoryResp = `
{
	"name": "Test Story",
	"user_id": 1,
	"description": null,
	"keep_events_for": 604800,
	"disabled": false,
	"priority": false,
	"send_to_story_enabled": false,
	"send_to_story_access_source": "OFF",
	"send_to_story_access": null,
	"shared_team_slugs": [],
	"entry_agent_id": null,
	"exit_agents": [],
	"send_to_story_skill_use_requires_confirmation": true,
	"team_id": 1,
	"tags": [],
	"guid": "a72744c526e7d3e5b608f130a583c98b",
	"slug": "new_story",
	"created_at": "2025-02-03T00:00:00Z",
	"updated_at": "2025-02-03T00:00:00Z",
	"edited_at": "2025-02-03T00:00:00Z",
	"mode": "LIVE",
	"id": 1,
	"folder_id": null,
	"published": true,
	"change_control_enabled": false,
	"locked": false,
	"owners": []
}`
)

func TestCreateStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusCreated, []byte(testStoryResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	draftStory := tines.Story{TeamID: 1}

	story, err := cli.CreateStory(ctx, &draftStory)

	assert.Nil(err, "the Tines client should create a story successfully")
	assert.Equal(1, story.TeamID, "the Tines client should create a story in the correct team")

}

func TestGetStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, []byte(testStoryResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	story, err := cli.GetStory(ctx, 1)

	assert.Nil(err, "the Tines client should create a story successfully")
	assert.Equal(story.Name, "New Story", "the Tines client should retrieve and parse the story successfully")

}

func TestUpdateStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, []byte(testUpdateStoryResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	update := tines.Story{Description: "Description"}

	story, err := cli.UpdateStory(ctx, 1, &update)

	assert.Nil(err, "the Tines client should update a story successfully")
	assert.Equal("Description", story.Description, "the Tines client should update the story successfully")

}

func TestListStories(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, []byte(testListStoriesResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	lf := tines.NewListFilter(
		tines.WithStoryFilter(tines.FilterPublished),
		tines.WithStoryOrder(tines.OrderByNameAsc),
	)

	storyList := cli.ListStories(ctx, lf)

	for s, err := range storyList {
		assert.Nil(err, "the list of stories should be iterable")
		assert.Equal("Test Story", s.Name, "the story name should be retrieved successfully")
	}
}

func TestDeleteStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusNoContent, nil)
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	err = cli.DeleteStory(ctx, 1)

	assert.Nil(err, "the Tines client should delete the story successfully")
}

func TestBatchDeleteStories(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusNoContent, nil)
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	err = cli.BatchDeleteStories(ctx, []int{1, 2})

	assert.Nil(err, "the Tines client should delete the batch of stories successfully")
}

func TestExportStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, []byte(testExportStoryResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	res, err := cli.ExportStory(ctx, 1, false)
	assert.Nil(err, "the Tines client should export a story successfully")
	assert.Equal("Test Story", res["name"], "the exported story should be valid JSON")

}

func TestImportStory(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, []byte(testImportStoryResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	var data map[string]interface{}
	file, err := os.ReadFile("./testdata/test-import.json")
	assert.Nil(err, "the test file for import should be read successfully")
	if err != nil {
		return
	}

	err = json.Unmarshal(file, &data)
	assert.Nil(err, "the imported story should be valid JSON")
	if err != nil {
		return
	}

	name, ok := data["name"].(string)
	assert.True(ok, "the imported story name should be a string")
	if !ok {
		return
	}

	sir := tines.StoryImportRequest{
		Data:    data,
		NewName: name,
		TeamID:  1,
		Mode:    tines.StoryModeReplace,
	}

	ctx := context.Background()

	res, err := cli.ImportStory(ctx, &sir)

	assert.Nil(err, "the story should import successfully")
	assert.Equal(name, res.Name, "the imported story name should match the JSON export file")

}
