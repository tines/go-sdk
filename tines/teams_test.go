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

func createTestTeamServer(handler http.HandlerFunc) (*httptest.Server, *tines.Client) {
	server := httptest.NewServer(handler)
	client, _ := tines.NewClient(
		tines.SetTenantUrl(server.URL),
		tines.SetApiKey("test-api-key"),
	)
	return server, client
}

const (
	// Sanitized API response value for team creation/retrieval
	testTeamResp = `{
		"id": 1,
		"name": "Security team",
		"groups": [
			{
				"id": 3,
				"name": "Tier 1"
			}
		]
	}`

	// Sanitized API response value for team listing
	testTeamListResp = `{
		"teams": [
			{
				"id": 1,
				"name": "Security team"
			},
			{
				"id": 2,
				"name": "Engineering team"
			}
		],
		"meta": {
			"pages": 1,
			"current_page": "1",
			"previous_page": null,
			"next_page": "",
			"next_page_number": 0,
			"per_page": 20,
			"count": 2
		}
	}`

	// Sanitized API response value for team member listing
	testTeamMemberListResp = `{
		"members": [
			{
				"id": 1,
				"email": "user1@example.com",
				"name": "User One",
				"role": "admin"
			},
			{
				"id": 2,
				"email": "user2@example.com",
				"name": "User Two",
				"role": "member"
			}
		],
		"meta": {
			"pages": 1,
			"current_page": "1",
			"previous_page": null,
			"next_page": "",
			"next_page_number": 0,
			"per_page": 20,
			"count": 2
		}
	}`
)

func TestCreateTeam(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var team tines.Team
		err := json.NewDecoder(r.Body).Decode(&team)
		assert.NoError(t, err)
		assert.Equal(t, "Security team", team.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testTeamResp))
	})
	defer server.Close()

	team := &tines.Team{
		Name: "Security team",
	}

	result, err := client.CreateTeam(context.Background(), team)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Security team", result.Name)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, 3, result.Groups[0].ID)
	assert.Equal(t, "Tier 1", result.Groups[0].Name)
}

func TestGetTeam(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testTeamResp))
	})
	defer server.Close()

	result, err := client.GetTeam(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Security team", result.Name)
}

func TestUpdateTeam(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var team tines.Team
		err := json.NewDecoder(r.Body).Decode(&team)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Security team", team.Name)

		// Return updated team
		updatedResp := `{
			"id": 1,
			"name": "Updated Security team",
			"groups": []
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(updatedResp))
	})
	defer server.Close()

	team := &tines.Team{
		Name: "Updated Security team",
	}

	result, err := client.UpdateTeam(context.Background(), 1, team)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Updated Security team", result.Name)
}

func TestListTeams(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testTeamListResp))
	})
	defer server.Close()

	filter := tines.NewListFilter()
	teams := make([]tines.Team, 0)

	for team, err := range client.ListTeams(context.Background(), filter) {
		assert.NoError(t, err)
		teams = append(teams, team)
	}

	assert.Len(t, teams, 2)
	assert.Equal(t, "Security team", teams[0].Name)
	assert.Equal(t, "Engineering team", teams[1].Name)
}

func TestListTeamsWithIncludePersonalTeams(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		// Verify that the include_personal_teams parameter is sent
		includePersonal := r.URL.Query().Get("include_personal_teams")
		assert.Equal(t, "true", includePersonal)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testTeamListResp))
	})
	defer server.Close()

	filter := tines.NewListFilter(tines.WithIncludePersonalTeams(true))
	teams := make([]tines.Team, 0)

	for team, err := range client.ListTeams(context.Background(), filter) {
		assert.NoError(t, err)
		teams = append(teams, team)
	}

	assert.Len(t, teams, 2)
	assert.Equal(t, "Security team", teams[0].Name)
	assert.Equal(t, "Engineering team", teams[1].Name)
}

func TestDeleteTeam(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteTeam(context.Background(), 1)
	assert.NoError(t, err)
}

func TestListTeamMembers(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testTeamMemberListResp))
	})
	defer server.Close()

	filter := tines.NewListFilter()
	members := make([]tines.TeamMember, 0)

	for member, err := range client.ListTeamMembers(context.Background(), 1, filter) {
		assert.NoError(t, err)
		members = append(members, member)
	}

	assert.Len(t, members, 2)
	assert.Equal(t, "user1@example.com", members[0].Email)
	assert.Equal(t, "User One", members[0].Name)
	assert.Equal(t, "admin", members[0].Role)
	assert.Equal(t, "user2@example.com", members[1].Email)
	assert.Equal(t, "User Two", members[1].Name)
	assert.Equal(t, "member", members[1].Role)
}

func TestInviteTeamMember(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var member tines.TeamMember
		err := json.NewDecoder(r.Body).Decode(&member)
		assert.NoError(t, err)
		assert.Equal(t, "newuser@example.com", member.Email)

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	member := &tines.TeamMember{
		Email: "newuser@example.com",
	}

	err := client.InviteTeamMember(context.Background(), 1, member)
	assert.NoError(t, err)
}

func TestRemoveTeamMember(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.RemoveTeamMember(context.Background(), 1, 2)
	assert.NoError(t, err)
}

func TestResendTeamMemberInvitation(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.ResendTeamMemberInvitation(context.Background(), 1, 2)
	assert.NoError(t, err)
}

// Error handling tests
func TestCreateTeamError(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": [{"message": "Name is required"}]}`))
	})
	defer server.Close()

	team := &tines.Team{} // Empty team should trigger validation error

	result, err := client.CreateTeam(context.Background(), team)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTeamNotFound(t *testing.T) {
	server, client := createTestTeamServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors": [{"message": "Team not found"}]}`))
	})
	defer server.Close()

	result, err := client.GetTeam(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, result)
}
