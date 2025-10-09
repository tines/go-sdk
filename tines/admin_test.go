package tines_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

const testCreateUserResp = `{
  "id": 276,
  "scim_uid": null,
  "email": "alice@tines.xyz",
  "created_at": "2019-11-03T09:57:49.537Z",
  "updated_at": "2019-11-03T09:57:49.537Z",
  "admin": false,
  "is_active": true,
  "first_name": "Alice",
  "last_name": "Smith",
  "last_seen": null,
  "tenant_permissions": ["AUDIT_LOG_READ"]
}`

const testGetUserResp = `{
  "id": 276,
  "scim_uid": null,
  "email": "alice@tines.xyz",
  "created_at": "2019-11-03T09:57:49.537Z",
  "updated_at": "2019-11-03T09:57:49.537Z",
  "admin": false,
  "is_active": true,
  "first_name": "Alice",
  "last_name": "Smith",
  "last_seen": "2025-01-15T10:30:00Z",
  "tenant_permissions": ["AUDIT_LOG_READ"]
}`

const testUpdateUserResp = `{
  "id": 276,
  "scim_uid": null,
  "email": "alice@tines.xyz",
  "created_at": "2019-11-03T09:57:49.537Z",
  "updated_at": "2025-01-15T12:00:00Z",
  "admin": true,
  "is_active": true,
  "first_name": "Alice",
  "last_name": "Jones",
  "last_seen": "2025-01-15T10:30:00Z",
  "tenant_permissions": ["AUDIT_LOG_READ", "TENANT_OWNER"]
}`

const testListUsersResp = `{
  "admin/users": [
    {
      "id": 2,
      "scim_uid": null,
      "email": "alice@tines.xyz",
      "created_at": "2019-07-08T14:16:11.967Z",
      "updated_at": "2019-10-31T21:40:52.682Z",
      "admin": false,
      "is_active": true,
      "first_name": "Alice",
      "last_name": "Smith",
      "last_seen": "2019-10-31T21:40:42.000Z",
      "tenant_permissions": ["AUDIT_LOG_READ"]
    },
    {
      "id": 3,
      "scim_uid": null,
      "email": "joe@tines.xyz",
      "created_at": "2019-09-11T14:46:37.880Z",
      "updated_at": "2019-09-27T20:45:46.527Z",
      "is_active": true,
      "first_name": "Joe",
      "last_name": "Bloggs",
      "last_seen": "2019-09-27T20:45:46.000Z",
      "tenant_permissions": []
    }
  ],
  "meta": {
    "current_page": "https://<<META.tenant.domain>>/api/v1/admin/users?page=1&per_page=20",
    "previous_page": null,
    "next_page": null,
    "next_page_number": null,
    "per_page": 20,
    "pages": 1,
    "count": 2
  }
}`

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusCreated, nil, []byte(testCreateUserResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	draft := tines.User{
		Id:                276,
		Email:             "alice@tines.xyz",
		FirstName:         "Alice",
		LastName:          "Smith",
		Admin:             lo.ToPtr(false),
		IsActive:          lo.ToPtr(true),
		TenantPermissions: []string{"AUDIT_LOG_READ"},
	}

	user, err := cli.CreateUser(context.Background(), &draft)

	assert.Nil(err, "the Tines client should create a user successfully")
	assert.Equal(draft.Id, user.Id, "the created user ID should match the response")
	assert.Equal(draft.Email, user.Email, "the created user email should match the response")
	assert.Equal(draft.FirstName, user.FirstName, "the created user first name should match the response")
	assert.Equal(draft.LastName, user.LastName, "the created user last name should match the response")
	assert.Equal(draft.Admin, user.Admin, "the created user admin status should match the response")
	assert.Equal(draft.IsActive, user.IsActive, "the created user active status should match the response")
	assert.Equal(draft.TenantPermissions, user.TenantPermissions, "the created user tenant permissions should match the response")
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testGetUserResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	user, err := cli.GetUser(context.Background(), 276)

	assert.Nil(err, "the Tines client should get a user successfully")
	assert.Equal(276, user.Id, "the user ID should match the response")
	assert.Equal("alice@tines.xyz", user.Email, "the user email should match the response")
	assert.Equal("Alice", user.FirstName, "the user first name should match the response")
	assert.Equal("Smith", user.LastName, "the user last name should match the response")
	assert.Equal(lo.ToPtr(false), user.Admin, "the user admin status should match the response")
	assert.Equal(lo.ToPtr(true), user.IsActive, "the user active status should match the response")
	assert.Equal("2025-01-15T10:30:00Z", user.LastSeen, "the user last seen should match the response")
	assert.Equal([]string{"AUDIT_LOG_READ"}, user.TenantPermissions, "the user tenant permissions should match the response")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testUpdateUserResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	updates := tines.User{
		Admin:             lo.ToPtr(true),
		LastName:          "Jones",
		TenantPermissions: []string{"AUDIT_LOG_READ", "TENANT_OWNER"},
	}

	user, err := cli.UpdateUser(context.Background(), 276, &updates)

	assert.Nil(err, "the Tines client should update a user successfully")
	assert.Equal(276, user.Id, "the updated user ID should match the response")
	assert.Equal("alice@tines.xyz", user.Email, "the updated user email should match the response")
	assert.Equal("Alice", user.FirstName, "the updated user first name should match the response")
	assert.Equal("Jones", user.LastName, "the updated user last name should match the response")
	assert.Equal(lo.ToPtr(true), user.Admin, "the updated user admin status should match the response")
	assert.Equal(lo.ToPtr(true), user.IsActive, "the updated user active status should match the response")
	assert.Equal([]string{"AUDIT_LOG_READ", "TENANT_OWNER"}, user.TenantPermissions, "the updated user tenant permissions should match the response")
}

func TestListUsers(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, []byte(testListUsersResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	var users []tines.User
	for u, err := range cli.ListUsers(context.Background(), tines.NewListFilter()) {
		assert.Nil(err, "no error should occur while iterating users")
		users = append(users, u)
	}

	assert.Equal(2, len(users), "should return 2 users")

	// Verify first user
	assert.Equal(2, users[0].Id, "first user ID should match")
	assert.Equal("alice@tines.xyz", users[0].Email, "first user email should match")
	assert.Equal("Alice", users[0].FirstName, "first user first name should match")
	assert.Equal("Smith", users[0].LastName, "first user last name should match")
	assert.Equal(lo.ToPtr(false), users[0].Admin, "first user admin status should match")
	assert.Equal(lo.ToPtr(true), users[0].IsActive, "first user active status should match")
	assert.Equal([]string{"AUDIT_LOG_READ"}, users[0].TenantPermissions, "first user permissions should match")

	// Verify second user
	assert.Equal(3, users[1].Id, "second user ID should match")
	assert.Equal("joe@tines.xyz", users[1].Email, "second user email should match")
	assert.Equal("Joe", users[1].FirstName, "second user first name should match")
	assert.Equal("Bloggs", users[1].LastName, "second user last name should match")
	assert.Nil(users[1].Admin, "second user admin status should be empty")
	assert.Equal(lo.ToPtr(true), users[1].IsActive, "second user active status should match")
	assert.Equal([]string{}, users[1].TenantPermissions, "second user permissions should match")
}

func TestDeleteUser(t *testing.T) {
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

	err = cli.DeleteUser(context.Background(), 276)

	assert.Nil(err, "the Tines client should delete a user successfully")
}
