package tines

import (
	"encoding/json"
	"strings"
	"time"
)

type ListFilter struct {
	TeamID      int    `json:"team_id,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Before      string `json:"before,omitempty"`
	After       string `json:"after,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
	OpName      string `json:"operation_name,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	Page        int    `json:"page,omitempty"`
	maxResults  int
}

// Filter results returned by a List endpoint (eg List Credentials, List Stories, etc).
//
// Example Usage:
//
//	lf := filter.NewListFilter(
//	  filter.WithTeamId(1),
//	)
func NewListFilter(opts ...func(*ListFilter)) ListFilter {
	lf := ListFilter{}

	for _, o := range opts {
		o(&lf)
	}

	return lf
}

// Limit results returned by a List endpoint to only the results that belong to a particular Team ID.
func WithTeamId(id int) func(*ListFilter) {
	return func(lf *ListFilter) {
		if id > 0 {
			lf.TeamID = id
		}
	}
}

// Limit results returned by a List endpoint to only the results that belong to a particular User ID.
func WithUserId(id int) func(*ListFilter) {
	return func(lf *ListFilter) {
		if id > 0 {
			lf.UserID = id
		}
	}
}

// Limit results returned by the List Folders endpoint to only folders that contain a certain type of content
// (eg Credentials, Resources, or Stories).
func WithContentType(ct string) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.ContentType = strings.ToUpper(ct)
	}
}

// Limit results returned by the List Audit Logs endpoint to only audit logs for a specific action or operation.
// The current list of logged operations is available at https://www.tines.com/api/audit-logs/. Logged operation
// names are case-sensitive.
func WithOperationName(op string) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.OpName = op
	}
}

// Limit results to those created before the specified ISO 8601 timestamp.
func WithResultsBefore(s string) func(*ListFilter) {
	return func(lf *ListFilter) {
		ts, err := time.Parse(time.RFC3339, s)
		if err != nil {
			lf.Before = ts.String()
		}
	}
}

// Limit results to those created after the specified ISO 8601 timestamp.
func WithResultsAfter(s string) func(*ListFilter) {
	return func(lf *ListFilter) {
		ts, err := time.Parse(time.RFC3339, s)
		if err != nil {
			lf.After = ts.String()
		}
	}
}

// Specify the number of results returned per page (minimum 20, maximum 500).
func WithResultsPerPage(i int) func(*ListFilter) {
	return func(lf *ListFilter) {
		switch {
		case i > 500:
			lf.PerPage = 500
		default:
			lf.PerPage = i
		}
	}
}

func WithResultsPageCursor(i int) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.Page = i
	}
}

func WithMaxResults(i int) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.maxResults = i
	}
}

func (l *ListFilter) AppendFilter(opt func(*ListFilter)) {
	opt(l)
}

func (l *ListFilter) MaxResults() int {
	return l.maxResults
}

func (l *ListFilter) ToParamMap() map[string]any {
	params := make(map[string]any)

	data, err := json.Marshal(l)
	if err != nil {
		params["error"] = err.Error()
		return params
	}

	err = json.Unmarshal(data, &params)
	if err != nil {
		params["error"] = err.Error()
	}

	return params
}
