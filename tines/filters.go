package tines

import (
	"encoding/json"
	"strings"
	"time"
)

// Enum for filtering a List of results.
type ResultFilter string

const (
	// Story List result filters.
	FilterApiEnabled ResultFilter = "API_ENABLED"
	FilterChangeCtrl ResultFilter = "CHANGE_CONTROL_ENABLED"
	FilterDisabled   ResultFilter = "DISABLED"
	FilterFavorite   ResultFilter = "FAVORITE"
	FilterHiPriority ResultFilter = "HIGH_PRIORITY"
	FilterLocked     ResultFilter = "LOCKED"
	FilterPublished  ResultFilter = "PUBLISHED"
	FilterStsEnabled ResultFilter = "SEND_TO_STORY_ENABLED"

	// Filter a list of Credentials to only those which are restricted.
	// This filter is only valid for lists of Credentials.
	FilterRestricted ResultFilter = "RESTRICTED"
	// Filter a list of Credentials to only those which are unrestricted.
	// This filter is only valid for lists of Credentials.
	FilterUnrestricted ResultFilter = "UNRESTRICTED"
	// Filter a list of Credentials to only those which are unused in actions.
	// This filter is only valid for lists of Credentials.
	FilterUnusedInActions ResultFilter = "UNUSED_IN_ACTIONS"
)

// Enum for ordering results returned by a List of Stories.
type StoryOrder string

const (
	OrderByActionCtAsc        StoryOrder = "ACTION_COUNT_ASC"
	OrderByActionCtDesc       StoryOrder = "ACTION_COUNT_DESC"
	OrderByNameAsc            StoryOrder = "NAME"
	OrderbyNameDesc           StoryOrder = "NAME_DESC"
	OrderByRecentlyEditedAsc  StoryOrder = "LEAST_RECENTLY_EDITED"
	OrderByRecentlyEditedDesc StoryOrder = "RECENTLY_EDITED"
)

type ListFilter struct {
	TeamID       int          `json:"team_id,omitempty"`
	FolderID     int          `json:"folder_id,omitempty"`
	ContentType  string       `json:"content_type,omitempty"`
	Before       string       `json:"before,omitempty"`
	After        string       `json:"after,omitempty"`
	UserID       int          `json:"user_id,omitempty"`
	OpName       string       `json:"operation_name,omitempty"`
	ResultFilter ResultFilter `json:"filter,omitempty"`
	StoryOrder   StoryOrder   `json:"order,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	PerPage      int          `json:"per_page,omitempty"`
	Page         int          `json:"page,omitempty"`
	maxResults   int
}

// Filter results returned by a List endpoint (eg List Credentials, List Stories, etc).
//
// Example Usage:
//
//	lf := tines.NewListFilter(
//	  tines.WithTeamId(1),
//	)
func NewListFilter(opts ...func(*ListFilter)) ListFilter {
	lf := ListFilter{
		maxResults: 100,
	}

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

// Limit results returned by a List endpoint to one of the enumerated ResultFilter types (eg API Enabled,
// Locked, etc.) Currently, the API allows only one filter to be applied at a time.
func WithResultFilter(f ResultFilter) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.ResultFilter = f
	}
}

// Specify the ordering of the results returned by the List Stories endpoint based on one of the enumerated fields
// (eg ordering by Story Name, Action Count, or Recently Edited, either ascending or descending.) Currently, the
// API only allows for specifying result order on one column at a time.
func WithStoryOrder(o StoryOrder) func(*ListFilter) {
	return func(lf *ListFilter) {
		lf.StoryOrder = o
	}
}

// Limit results to those created before the specified ISO 8601 timestamp.
func WithResultsBefore(s string) func(*ListFilter) {
	return func(lf *ListFilter) {
		ts, err := time.Parse(time.RFC3339, s)
		if err != nil {
			// Optionally fall back to a partially-compliant timestamp.
			ts2, err := time.Parse("2006-01-02", s)
			if err != nil {
				// If we still can't parse a value, default to current time.
				lf.Before = time.Now().Format(time.RFC3339)
			} else {
				lf.Before = ts2.Format(time.RFC3339)
			}
		} else {
			lf.Before = ts.Format(time.RFC3339)
		}
	}
}

// Limit results to those created after the specified ISO 8601 timestamp.
func WithResultsAfter(s string) func(*ListFilter) {
	return func(lf *ListFilter) {
		ts, err := time.Parse(time.RFC3339, s)
		if err != nil {
			// Optionally fall back to a partially-compliant timestamp.
			ts2, err := time.Parse("2006-01-02", s)
			if err != nil {
				// If we still can't parse a value, default to current time.
				lf.After = time.Now().Format(time.RFC3339)
			} else {
				lf.After = ts2.Format(time.RFC3339)
			}
		} else {
			lf.After = ts.Format(time.RFC3339)
		}
	}
}

// Specify the number of results returned per page (minimum 20, maximum 500).
func WithResultsPerPage(i int) func(*ListFilter) {
	return func(lf *ListFilter) {
		switch {
		case i > 500:
			lf.PerPage = 500
		case i < 20:
			lf.PerPage = 20
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
