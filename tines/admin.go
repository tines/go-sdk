package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type User struct {
	Id                int      `json:"id,omitempty"`
	ScimUid           string   `json:"scim_uid,omitempty"`
	Email             string   `json:"email,omitempty"`
	Admin             *bool    `json:"admin,omitempty"`
	TenantPermissions []string `json:"tenant_permissions,omitempty"`
	IsActive          *bool    `json:"is_active,omitempty"`
	CreatedAt         string   `json:"created_at,omitempty"`
	UpdatedAt         string   `json:"updated_at,omitempty"`
	FirstName         string   `json:"first_name,omitempty"`
	LastName          string   `json:"last_name,omitempty"`
	LastSeen          string   `json:"last_seen,omitempty"`
}

type UserList struct {
	Users []User        `json:"admin/users,omitempty"`
	Meta  paginate.Meta `json:"meta,omitempty"`
}

// CreateUser creates a new User.
func (c *Client) CreateUser(ctx context.Context, u *User) (*User, error) {
	resource := "/api/v1/admin/users"
	errs := Error{Type: ErrorTypeRequest}
	newUser := User{}

	if u.Email == "" {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "User Email must not be empty",
		})
	}

	if errs.HasErrors() {
		return nil, errs
	}

	req, err := json.Marshal(u)
	if err != nil {
		return &newUser, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errParseError,
					Details: err.Error(),
				},
			},
		}
	}

	body, err := c.doRequest(ctx, http.MethodPost, resource, nil, req)
	if err != nil {
		return &newUser, err
	}

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		return &newUser, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &newUser, nil
}

// GetUser retrieves a user by its ID.
func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	u := User{}
	resource := fmt.Sprintf("/api/v1/admin/users/%d", id)

	body, err := c.doRequest(ctx, http.MethodGet, resource, nil, nil)
	if err != nil {
		return &u, err
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return &u, err
	}

	return &u, nil
}

// UpdateUser updates a user by its ID.
func (c *Client) UpdateUser(ctx context.Context, id int, u *User) (*User, error) {
	resource := fmt.Sprintf("/api/v1/admin/users/%d", id)
	updatedUser := User{}

	req, err := json.Marshal(u)
	if err != nil {
		return &updatedUser, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errParseError,
					Details: err.Error(),
				},
			},
		}
	}

	body, err := c.doRequest(ctx, http.MethodPut, resource, nil, req)
	if err != nil {
		return &updatedUser, err
	}

	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return &updatedUser, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &updatedUser, nil
}

// ListUsers yields an iterator that returns individual Users, optionally filtered by UserFilter.
//
// If no filters are specified, ListUsers() will recurse through all pages of results
// until no more are available.
//
// If `filters.WithMaxResults()` is set, this function will
// yield either the actual set of results or the specified maximum number of results,
// whichever is less.
//
// Example Usage:
//
//	for u, err := range ListUsers(ctx, filters.WithMaxResults(10)) {
//		if err != nil {
//			...
//		}
//		fmt.Println(u.Email)
//	}
func (c *Client) ListUsers(ctx context.Context, f ListFilter) iter.Seq2[User, error] {
	var userList, resultList UserList
	resource := "/api/v1/admin/users"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	return func(yield func(User, error) bool) {
		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(User{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(User{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.Users {
				userList.Users = append(userList.Users, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			resultList = UserList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range userList.Users {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// DeleteUser deletes a user by its ID.
func (c *Client) DeleteUser(ctx context.Context, id int) error {
	resource := fmt.Sprintf("/api/v1/admin/users/%d", id)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
