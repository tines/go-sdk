package tines_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

func TestCreateStoryVersion(t *testing.T) {
	assert := assert.New(t)

	storyVersion := tines.StoryVersion{
		ID:          885,
		Name:        "API created version",
		Description: "",
		Timestamp:   "2023-11-13T11:16:05Z",
	}

	respBody, _ := json.Marshal(storyVersion)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("POST", r.Method)
		assert.Equal("/api/v1/stories/123/versions", r.URL.Path)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("Bearer foo", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}))
	defer ts.Close()

	client, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)
	assert.Nil(err)

	request := &tines.StoryVersionCreateRequest{
		Name: "API created version",
	}

	result, err := client.CreateStoryVersion(context.Background(), 123, request)

	assert.Nil(err)
	assert.Equal(885, result.ID)
	assert.Equal("API created version", result.Name)
	assert.Equal("", result.Description)
	assert.Equal("2023-11-13T11:16:05Z", result.Timestamp)
}
