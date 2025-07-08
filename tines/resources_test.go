package tines_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

const (
	// Sanitized API request as of 2025-06-23.
	testCreateTextResourceReq = `
{
    "name": "Test",
    "value": "value",
    "team_id": 1
}`
	// Sanitized API response as of 2025-06-23.
	testCreateTextResourceResp = `
{
    "id": 1,
    "name": "Test",
    "value": "value",
    "team_id": 1,
    "folder_id": null,
    "user_id": 1,
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "test_resource_enabled": false,
    "referencing_action_ids": []
}`
	// Sanitized API request as of 2025-06-23.
	testCreateJsonResourceReq = `
{
    "name": "Test 2",
    "value": {"foo":"bar"},
    "team_id": 1
}`
	// Sanitized API response as of 2025-06-23.
	testCreateJsonResourceResp = `
{
    "id": 1,
    "name": "Test 2",
    "value": "{\"foo\":\"bar\"}",
    "team_id": 1,
    "folder_id": null,
    "user_id": 1,
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test_2",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "test_resource_enabled": false,
    "referencing_action_ids": []
}`
	// Santiized API response as of 2025-06-23.
	testUpdateResourceResp = `
{
    "id": 1,
    "name": "Test",
    "value": "value",
    "team_id": 1,
    "folder_id": null,
    "user_id": 1,
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "test_resource_enabled": false,
    "referencing_action_ids": []
}`
	// Santiized API response as of 2025-06-23.
	testUpdateEmptyResourceResp = `
{
    "id": 1,
    "name": "Test",
    "value": "",
    "team_id": 1,
    "folder_id": null,
    "user_id": 1,
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "test_resource_enabled": false,
    "referencing_action_ids": []
}`
	// Sanitized API response as of 2025-06-23.
	testListResourcesResp = `
{
    "global_resources": [
		{
			"id": 1,
			"name": "Test",
			"value": "value",
			"team_id": 1,
			"folder_id": null,
			"user_id": 1,
			"read_access": "TEAM",
			"shared_team_slugs": [],
			"slug": "test",
			"created_at": "2025-01-01T00:00:00Z",
			"updated_at": "2025-01-01T00:00:00Z",
			"description": "",
			"test_resource_enabled": false,
			"referencing_action_ids": []
		}
    ],
    "meta": {
        "current_page": "https://example.tines.com/api/v1/global_resources?per_page=20&page=1",
        "previous_page": null,
        "next_page": null,
        "next_page_number": null,
        "per_page": 20,
        "pages": 1,
        "count": 1
    }
}`
	testAppendStringReq = `
{
	"value": "foo"
}
`
	testAppendStringResp = `valuefoo`
	testAppendArrayReq   = `
{
    "value": [2]
}`

	testAppendArrayResp = `["one", 2]`
)

func TestCreateResource(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		req     string
		resp    string
		payload tines.Resource
	}{
		{"TextCredential", testCreateTextResourceReq, testCreateTextResourceResp, tines.Resource{Name: "Test", Value: "value", TeamId: 1}},
		{"JsonCredential", testCreateJsonResourceReq, testCreateJsonResourceResp, tines.Resource{Name: "Test 2", Value: map[string]any{"foo": "bar"}, TeamId: 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := createTestServer(assert, http.StatusOK, []byte(test.req), []byte(test.resp))
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

			res, err := cli.CreateResource(ctx, &test.payload)

			assert.Nil(err, "the resource should be created without errors")
			assert.IsType(&tines.Resource{}, res, "the response should be the expected type")
			assert.Equal(test.payload.Name, res.Name, "the created credential name should match the request")
		})
	}

}

func TestGetResource(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testCreateTextCredResp))
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

	res, err := cli.GetResource(ctx, 1)
	assert.Nil(err, "the resource should be created without errors")
	assert.IsType(&tines.Resource{}, res, "the response should be the expected type")
	assert.Equal(1, res.Id, "the retrieved credential should match the request")
}

func TestUpdateResource(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		req     string
		resp    string
		payload tines.Resource
	}{
		{"WithValue", "", testUpdateResourceResp, tines.Resource{Id: 1, Name: "Test", Value: "value", TeamId: 1}},
		{"EmptyValue", "", testUpdateEmptyResourceResp, tines.Resource{Id: 1, Name: "Test 2", Value: "", TeamId: 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := createTestServer(assert, http.StatusOK, nil, []byte(test.resp))
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

			res, err := cli.UpdateResource(ctx, 1, &test.payload)
			assert.Nil(err, "the resource should be retrieved without errors")
			assert.IsType(&tines.Resource{}, res, "the response should be the expected type")
			assert.Equal(test.payload.Value, res.Value, "the retrieved credential should match the request")

		})
	}
}

func TestListResources(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testListResourcesResp))
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

	lf := tines.NewListFilter()

	resList := cli.ListResources(ctx, lf)

	for c, err := range resList {
		assert.Nil(err, "the list of resources should be iterable")
		assert.Equal("Test", c.Name, "the resource name should be retrieved successfully")
	}
}

func TestAppendResourceElement(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		req     string
		resp    string
		prevVal any
		payload tines.ResourceElement
	}{
		{"TestAppendStringToString", testAppendStringReq, testAppendStringResp, "value", tines.ResourceElement{Value: "foo"}},
		{"TestAppendArrayToArray", testAppendArrayReq, testAppendArrayResp, []any{"one"}, tines.ResourceElement{Value: []any{2}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := createTestServer(assert, http.StatusOK, []byte(test.req), []byte(test.resp))
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

			res, err := cli.AppendResourceElement(ctx, 1, &test.payload)
			assert.Nil(err, "the resource should be created without errors")
			assert.IsType("", res, "the response should be the expected type")
		})
	}
}
