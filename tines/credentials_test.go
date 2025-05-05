package tines_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

const (
	// Sanitized API request as of 2025-04-09.
	testCreateAwsCredReq = `
{
    "name": "Test",
    "team_id": 1,
    "mode": "AWS",
    "aws_authentication_type": "KEY",
    "aws_access_key": "foo",
    "aws_secret_key": "bar"
}`

	// Sanitized API response as of 2025-04-09.
	testCreateAwsCredResp = `
{
    "id": 1,
    "name": "Test",
    "team_id": 1,
    "folder_id": null,
    "mode": "AWS",
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "aws_assumed_role_external_id": "2f1a14a8-35c8-473e-ad4a-8bcc64e02b6b",
    "aws_authentication_type": "KEY",
    "allowed_hosts": [],
    "metadata": {},
    "restriction_type": "RESTRICTED",
    "test_credential_enabled": false,
    "test_credential": null
}`

	// Sanitized API request as of 2025-04-09.
	testCreateHttpCredReq = `
{
    "name": "Test",
    "team_id": 1,
    "mode": "HTTP_REQUEST_AGENT",
    "http_request_options": {
        "url": "https://example.com",
        "content_type": "json",
        "method": "post",
        "payload": {
            "key": "<<secret>>"
        }
    },
    "http_request_location_of_token": "=response.body.value",
    "http_request_secret": "example",
    "http_request_ttl": 60
}`
	// Sanitized APIU response as of 2025-04-09.
	testCreateHttpCredResp = `
{
    "id": 1,
    "name": "Test",
    "team_id": 1,
    "folder_id": null,
    "mode": "HTTP_REQUEST_AGENT",
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "aws_assumed_role_external_id": null,
    "aws_authentication_type": null,
    "allowed_hosts": [],
    "metadata": {},
    "restriction_type": "RESTRICTED",
    "test_credential_enabled": false,
    "test_credential": null
}`

	// Sanitized API request as of 2025-04-09.
	testCreateTextCredReq = `
{
    "name": "Test",
    "team_id": 1,
    "mode": "TEXT",
    "value": "value"
}`

	// Sanitized API response as of 2025-04-09.
	testCreateTextCredResp = `
{
    "id": 1,
    "name": "Test",
    "team_id": 1,
    "folder_id": null,
    "mode": "TEXT",
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "aws_assumed_role_external_id": null,
    "aws_authentication_type": null,
    "allowed_hosts": [],
    "metadata": {},
    "restriction_type": "RESTRICTED",
    "test_credential_enabled": false,
    "test_credential": null
}`

	// Sanitized API response as of 2025-04-09.
	testGetCredResp = `
{
    "id": 1,
    "name": "Test",
    "team_id": 1,
    "folder_id": null,
    "mode": "HTTP_REQUEST_AGENT",
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "description": "",
    "aws_assumed_role_external_id": null,
    "aws_authentication_type": null,
    "allowed_hosts": [],
    "metadata": {},
    "restriction_type": "RESTRICTED",
    "test_credential_enabled": false
}`

	// Sanitized API response as of 2025-04-09.
	testUpdateCredResp = `
{
    "id": 1,
    "name": "Test",
    "team_id": 1,
    "folder_id": null,
    "mode": "TEXT",
    "read_access": "TEAM",
    "shared_team_slugs": [],
    "slug": "test",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-02T00:00:00Z",
    "description": "",
    "aws_assumed_role_external_id": null,
    "aws_authentication_type": null,
    "allowed_hosts": [],
    "metadata": {},
    "restriction_type": "RESTRICTED",
    "test_credential_enabled": false,
    "test_credential": null
}`

	// Sanitized API response as of 2025-04-09.
	testListCredsResp = `
{
    "user_credentials": [
        {
            "id": 1,
            "name": "example test credential",
            "team_id": 1,
            "folder_id": null,
            "mode": "TEXT",
            "read_access": "TEAM",
            "shared_team_slugs": [],
            "slug": "example_text_credential",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z",
            "description": "",
            "aws_assumed_role_external_id": null,
            "aws_authentication_type": null,
            "allowed_hosts": [],
            "metadata": {},
            "restriction_type": "RESTRICTED",
            "test_credential_enabled": false,
            "test_credential": null
        }
    ],
    "meta": {
        "current_page": "https://example.tines.com/api/v1/user_credentials?per_page=1&page=1&filter=RESTRICTED",
        "previous_page": null,
        "next_page": null,
        "next_page_number": 1,
        "per_page": 1,
        "pages": 1,
        "count": 1
    }
}`
)

func TestCreateCredential(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name   string
		req    string
		resp   string
		mode   tines.CredentialType
		values tines.CredentialPayload
	}{
		{"AwsCredential", testCreateAwsCredReq, testCreateAwsCredResp, tines.CredentialTypeAws, tines.CredentialPayload{AwsAuthType: "KEY", AwsAccessKey: "foo", AwsSecretKey: "bar"}},
		{"HttpCredential", testCreateHttpCredReq, testCreateHttpCredResp, tines.CredentialTypeHttp, tines.CredentialPayload{HttpReqOpts: map[string]any{"url": "https://example.com", "content_type": "json", "method": "post", "payload": map[string]any{"key": "<<secret>>"}}, HttpReqTtl: 60, HttpReqSecret: "example", HttpReqTokenLoc: "=response.body.value"}},
		{"TextCredential", testCreateTextCredReq, testCreateTextCredResp, tines.CredentialTypeText, tines.CredentialPayload{TextValue: "value"}},
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

			cred := tines.Credential{
				Name:              "Test",
				TeamId:            1,
				Mode:              test.mode,
				CredentialPayload: test.values,
			}

			res, err := cli.CreateCredential(ctx, &cred)

			assert.Nil(err, "the credential should be created without errors")
			assert.IsType(&tines.Credential{}, res, "the response should be the expected type")
			assert.Equal(cred.Name, res.Name, "the created credential name should match the request")

		})
	}
}

func TestGetCredential(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testGetCredResp))
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

	cred, err := cli.GetCredential(ctx, 1)
	assert.Nil(err, "the credential should be retrieved without errors")
	assert.IsType(&tines.Credential{}, cred, "the response should be the expected type")
	assert.Equal("Test", cred.Name, "the created credential name should match the request")

}

func TestUpdateCredential(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testUpdateCredResp))
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

	update := tines.Credential{
		Id:   1,
		Mode: tines.CredentialTypeText,
		CredentialPayload: tines.CredentialPayload{
			TextValue: "new",
		},
	}

	cred, err := cli.UpdateCredential(ctx, update.Id, &update)

	assert.Nil(err, "the credential should be retrieved without errors")
	assert.IsType(&tines.Credential{}, cred, "the response should be the expected type")

	created, err := time.Parse(time.RFC3339, cred.CreatedAt)
	assert.Nil(err, "the credential creation timestamp should be valid")

	updated, err := time.Parse(time.RFC3339, cred.UpdatedAt)
	assert.Nil(err, "the credential update timestamp should be valid")

	assert.Greater(updated, created, "the credential update time should be more recent than the creation time")

}

func TestListCredentials(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testListCredsResp))
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
		tines.WithResultFilter(tines.FilterRestricted),
	)

	credList := cli.ListCredentials(ctx, lf)

	for c, err := range credList {
		assert.Nil(err, "the list of stories should be iterable")
		assert.Equal("example test credential", c.Name, "the credential name should be retrieved successfully")
	}

}

func TestDeleteCredential(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusNoContent, nil, nil)
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

	err = cli.DeleteCredential(ctx, 1)

	assert.Nil(err, "the Tines client should delete the credential successfully")

}
