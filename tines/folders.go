package tines

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type Folder struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	TeamID      int    `json:"team_id,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size,omitempty"`
}

type FolderList struct {
	Folders []Folder      `json:"folders,omitempty"`
	Meta    paginate.Meta `json:"meta,omitempty"`
}

// Create a new Folder. Name, TeamID, and ContentType are required parameters.
func (c *Client) CreateFolder(ctx context.Context, f *Folder) (*Folder, error) {
	resource := "/api/v1/folders"
	errs := Error{Type: ErrorTypeRequest}

	if f.Name == "" {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "Folder Name must not be empty",
		})
	}

	if f.ContentType == "" || (f.ContentType != "CREDENTIAL" && f.ContentType != "RESOURCE" && f.ContentType != "STORY") {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "Folder Content Type must be one of \"CREDENTIAL\", \"RESOURCE\", or \"STORY\"",
		})
	}

	if f.TeamID == 0 {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: "Folder Team ID must not be empty",
		})
	}

	if errs.HasErrors() {
		return nil, errs
	}

	newFolder := Folder{
		Name:        f.Name,
		ContentType: f.ContentType,
		TeamID:      f.TeamID,
	}

	req, err := json.Marshal(&newFolder)
	if err != nil {
		return &newFolder, Error{
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
		return &newFolder, err
	}

	err = json.Unmarshal(body, &newFolder)
	if err != nil {
		return &newFolder, Error{
			Type: ErrorTypeServer,
			Errors: []ErrorMessage{
				{
					Message: errUnmarshalError,
					Details: err.Error(),
				},
			},
		}
	}

	return &newFolder, nil
}

// Get a Folder by unique ID.
func (c *Client) GetFolder(ctx context.Context, id int) (*Folder, error) {
	f := Folder{}
	resource := fmt.Sprintf("/api/v1/folders/%d", id)

	body, err := c.doRequest(ctx, http.MethodGet, resource, nil, nil)
	if err != nil {
		return &f, err
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		return &f, err
	}

	return &f, nil
}

// Update a Folder by unique ID. The only attribute that can be updated is the folder name.
// To change any other folder attribute, the resource must be deleted and recreated.
func (c *Client) UpdateFolder(ctx context.Context, id int, name string) (*Folder, error) {
	f := Folder{}
	resource := fmt.Sprintf("/api/v1/folders/%d", id)
	var params = make(map[string]any)

	params["name"] = name

	body, err := c.doRequest(ctx, http.MethodPut, resource, params, nil)
	if err != nil {
		return &f, err
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		return &f, err
	}
	return &f, nil
}

// Yields an iterator that returns individual Folders, optionally filtered by Team ID and/or
// by folder content type. If no other filters are specified, ListFolders() will recurse
// through all pages of results until no more are available. If `filters.WithMaxResults()` is
// set, this function will yield either the actual set of results or the specified maximum
// number of results, whichever is less.
//
// Example Usage:
//
//	for f, err := range ListFolders(ctx, filters.WithMaxResults(10)) {
//		if err != nil {
//			...
//		}
//		fmt.Println(f.Name)
//	}
func (c *Client) ListFolders(ctx context.Context, f ListFilter) iter.Seq2[Folder, error] {
	var folderList, resultList FolderList
	resource := "/api/v1/folders"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	c.logger.Debug(fmt.Sprintf("max results requested: %d", page.TotalRequested))

	return func(yield func(Folder, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(Folder{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(Folder{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.Folders {
				folderList.Folders = append(folderList.Folders, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = FolderList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range folderList.Folders {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// Delete a folder by unique ID.
func (c *Client) DeleteFolder(ctx context.Context, id int) error {
	resource := fmt.Sprintf("/api/v1/folders/%d", id)

	_, err := c.doRequest(ctx, http.MethodDelete, resource, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
