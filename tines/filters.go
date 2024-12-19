package tines

import (
	"encoding/json"
	"strings"
)

type ListFilter struct {
	TeamId      int    `json:"team_id,omitempty"`
	ContentType string `json:"content_type,omitempty"`
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
			lf.TeamId = id
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
