package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

// Team represents a Tines team.
type Team struct {
	ID     int         `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Groups []TeamGroup `json:"groups,omitempty"`
}

// TeamGroup represents a team group (sub-team within a team).
type TeamGroup struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TeamList represents paginated team results.
type TeamList struct {
	Teams []Team        `json:"teams,omitempty"`
	Meta  paginate.Meta `json:"meta,omitempty"`
}

// TeamMember represents a team member - used for both responses and invite requests.
type TeamMember struct {
	ID    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Role  string `json:"role,omitempty"`
}

// TeamMemberList represents paginated team member results.
type TeamMemberList struct {
	Members []TeamMember  `json:"members,omitempty"`
	Meta    paginate.Meta `json:"meta,omitempty"`
}

// CreateTeam creates a new team with the specified name.
func (c *Client) CreateTeam(ctx context.Context, team *Team) (*Team, error) {
	newTeam := Team{}

	req, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, "POST", "/api/v1/teams", nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &newTeam)
	if err != nil {
		return nil, err
	}

	return &newTeam, nil
}

// GetTeam retrieves details for a specific team by ID.
func (c *Client) GetTeam(ctx context.Context, id int) (*Team, error) {
	resource := fmt.Sprintf("/api/v1/teams/%d", id)
	team := Team{}

	body, err := c.doRequest(ctx, "GET", resource, nil, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &team)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// UpdateTeam updates an existing team's information.
func (c *Client) UpdateTeam(ctx context.Context, id int, team *Team) (*Team, error) {
	updatedTeam := Team{}
	resource := fmt.Sprintf("/api/v1/teams/%d", id)

	req, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, "PUT", resource, nil, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &updatedTeam)
	if err != nil {
		return nil, err
	}

	return &updatedTeam, nil
}

// ListTeams yields an iterator that returns individual Teams, optionally filtered.
// If no other filters are specified, ListTeams() will recurse through all pages of
// results until no more are available. If `filters.WithMaxResults()` is set, this
// function will yield either the actual set of results or the specified maximum
// number of results, whichever is less.
//
// Example Usage:
//
//	for t, err := range client.ListTeams(ctx, tines.NewListFilter()) {
//		if err != nil {
//			...
//		}
//		fmt.Println(t.Name)
//	}
func (c *Client) ListTeams(ctx context.Context, f ListFilter) iter.Seq2[Team, error] {
	var teamList, resultList TeamList
	resource := "/api/v1/teams"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	c.logger.Debug(fmt.Sprintf("max results requested: %d", page.TotalRequested))

	return func(yield func(Team, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(Team{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(Team{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.Teams {
				teamList.Teams = append(teamList.Teams, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = TeamList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range teamList.Teams {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// DeleteTeam deletes a team by ID.
func (c *Client) DeleteTeam(ctx context.Context, id int) error {
	resource := fmt.Sprintf("/api/v1/teams/%d", id)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)

	return err
}

// ListTeamMembers yields an iterator that returns individual TeamMembers for a specific team,
// optionally filtered. If no other filters are specified, ListTeamMembers() will recurse
// through all pages of results until no more are available. If `filters.WithMaxResults()` is
// set, this function will yield either the actual set of results or the specified maximum
// number of results, whichever is less.
//
// Example Usage:
//
//	for m, err := range client.ListTeamMembers(ctx, teamID, tines.NewListFilter()) {
//		if err != nil {
//			...
//		}
//		fmt.Println(m.Email)
//	}
func (c *Client) ListTeamMembers(ctx context.Context, teamID int, f ListFilter) iter.Seq2[TeamMember, error] {
	var memberList, resultList TeamMemberList
	resource := fmt.Sprintf("/api/v1/teams/%d/members", teamID)
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	c.logger.Debug(fmt.Sprintf("max results requested: %d", page.TotalRequested))

	return func(yield func(TeamMember, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(TeamMember{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(TeamMember{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.Members {
				memberList.Members = append(memberList.Members, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = TeamMemberList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range memberList.Members {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// InviteTeamMember invites a new member to join a team by email.
func (c *Client) InviteTeamMember(ctx context.Context, teamID int, member *TeamMember) error {
	resource := fmt.Sprintf("/api/v1/teams/%d/invite", teamID)

	req, err := json.Marshal(member)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, "POST", resource, nil, req)
	return err
}

// RemoveTeamMember removes a member from a team.
func (c *Client) RemoveTeamMember(ctx context.Context, teamID, memberID int) error {
	resource := fmt.Sprintf("/api/v1/teams/%d/members/%d", teamID, memberID)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)
	return err
}

// ResendTeamMemberInvitation resends an invitation to a team member.
func (c *Client) ResendTeamMemberInvitation(ctx context.Context, teamID, memberID int) error {
	resource := fmt.Sprintf("/api/v1/teams/%d/members/%d/resend", teamID, memberID)

	_, err := c.doRequest(ctx, "POST", resource, nil, nil)
	return err
}
