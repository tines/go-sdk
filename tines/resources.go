package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type Resource struct {
	Id             int      `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Value          any      `json:"value,omitempty"`
	TeamId         int      `json:"team_id,omitempty"`
	FolderId       int      `json:"folder_id,omitempty"`
	UserId         int      `json:"user_id,omitempty"`
	ReadAccess     string   `json:"read_access,omitempty"`
	SharedTeams    []string `json:"shared_team_slugs,omitempty"`
	Slug           string   `json:"slug,omitempty"`
	Description    string   `json:"description,omitempty"`
	TestResEnabled bool     `json:"test_resource_enabled,omitempty"`
	TestResource   any      `json:"test_resource,omitempty"`
	IsTest         bool     `json:"is_test,omitempty"`
	LiveResourceId int      `json:"live_resource_id,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
	RefActions     []int    `json:"referencing_action_ids,omitempty"`
}

type ResourceElement struct {
	ResourceId int    `json:"resource_id,omitempty"`
	Key        string `json:"key,omitempty"`
	Index      int    `json:"index,omitempty"`
	Value      any    `json:"value,omitempty"`
	IsTest     bool   `json:"is_test,omitempty"`
}

type ResourceList struct {
	GlobalResources []Resource    `json:"global_resources,omitempty"`
	Meta            paginate.Meta `json:"meta,omitempty"`
}

func (c *Client) CreateResource(ctx context.Context, r *Resource) (*Resource, error) {
	resource := "/api/v1/global_resources"
	newRes := Resource{}

	req, err := json.Marshal(r)
	if err != nil {
		return &newRes, Error{
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
		return &newRes, err
	}

	err = json.Unmarshal(body, &newRes)
	if err != nil {
		return &newRes, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &newRes, nil
}

func (c *Client) GetResource(ctx context.Context, id int) (*Resource, error) {
	resource := fmt.Sprintf("/api/v1/global_resources/%d", id)
	res := Resource{}

	body, err := c.doRequest(ctx, http.MethodGet, resource, nil, nil)
	if err != nil {
		return &res, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return &res, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &res, nil
}

func (c *Client) UpdateResource(ctx context.Context, id int, r *Resource) (*Resource, error) {
	resource := fmt.Sprintf("/api/v1/global_resources/%d", id)
	errs := Error{Type: ErrorTypeRequest}
	updatedRes := Resource{}

	if r.Id == 0 {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "You must specify the Resource ID to update.",
		})
	}

	if errs.HasErrors() {
		return nil, errs
	}

	req, err := json.Marshal(r)
	if err != nil {
		return &updatedRes, Error{
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
		return &updatedRes, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	err = json.Unmarshal(body, &updatedRes)
	if err != nil {
		return &updatedRes, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &updatedRes, nil

}

func (c *Client) ListResources(ctx context.Context, f ListFilter) iter.Seq2[Resource, error] {
	var resourceList, resultList ResourceList
	resource := "/api/v1/global_resources"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}
	return func(yield func(Resource, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(Resource{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(Resource{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.GlobalResources {
				resourceList.GlobalResources = append(resourceList.GlobalResources, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = ResourceList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range resourceList.GlobalResources {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func (c *Client) DeleteResource(ctx context.Context, id int) error {
	resource := fmt.Sprintf("/api/v1/global_resources/%d", id)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)
	if err != nil {
		return err
	}

	return nil

}

// Currently, it's only possible to append a string to a string, or an array to an array. This API
// endpoint currently returns the updated string value if the underlying Resource is a string, or an
// array if the underlying Resource is an array. To make the results consistent with the stringified
// Resource value returned by GetResource(), we automatically stringify array results.
func (c *Client) AppendResourceElement(ctx context.Context, id int, e *ResourceElement) (string, error) {
	resource := fmt.Sprintf("/api/v1/global_resources/%d/append", id)

	req, err := json.Marshal(e)
	if err != nil {
		return "", Error{
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
		return "", Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	return string(body), nil
}

// If a Resource value is a non-string JSON type, elements can be removed by key (if the resource is an object)
// or by index (if the resource is an array). Currently, the result of this API operation is returned as a string,
// which is not consistent with the results of the API operations for appending or replacing Resource elements.
func (c *Client) RemoveResourceElement(ctx context.Context, id int, e *ResourceElement) (string, error) {
	resource := fmt.Sprintf("/api/v1/global_resources/%d/remove", id)
	errs := Error{Type: ErrorTypeRequest}

	if e.ResourceId == 0 {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "You must specify the Resource ID to update.",
		})
	}

	if errs.HasErrors() {
		return "", errs
	}

	req, err := json.Marshal(e)
	if err != nil {
		return "", Error{
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
		return "", Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	return string(body), nil
}

func (c *Client) ReplaceResourceElement(ctx context.Context, id int, e *ResourceElement) (*ResourceElement, error) {
	resource := fmt.Sprintf("/api/v1/global_resources/%d/replace", id)
	errs := Error{Type: ErrorTypeRequest}
	updatedRes := ResourceElement{}

	if e.ResourceId == 0 {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "You must specify the Resource ID to update.",
		})
	}

	if errs.HasErrors() {
		return nil, errs
	}

	req, err := json.Marshal(e)
	if err != nil {
		return &updatedRes, Error{
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
		return &updatedRes, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	err = json.Unmarshal(body, &updatedRes)
	if err != nil {
		return &updatedRes, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &updatedRes, nil
}
