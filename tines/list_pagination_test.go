package tines_test

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tines/go-sdk/tines"
)

func listIDs[T any](seq iter.Seq2[T, error], id func(T) int) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		for value, err := range seq {
			if !yield(id(value), err) {
				return
			}
		}
	}
}

func TestListPagination(t *testing.T) {
	endpoints := []struct {
		name string
		path string
		key  string
		list func(*tines.Client, tines.ListFilter) iter.Seq2[int, error]
	}{
		{"stories", "/api/v1/stories", "stories", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListStories(context.Background(), f), func(v tines.Story) int { return v.ID })
		}},
		{"credentials", "/api/v1/user_credentials", "user_credentials", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListCredentials(context.Background(), f), func(v tines.Credential) int { return v.Id })
		}},
		{"resources", "/api/v1/global_resources", "global_resources", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListResources(context.Background(), f), func(v tines.Resource) int { return v.Id })
		}},
		{"folders", "/api/v1/folders", "folders", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListFolders(context.Background(), f), func(v tines.Folder) int { return v.Id })
		}},
		{"audit_logs", "/api/v1/audit_logs", "audit_logs", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListAuditLogs(context.Background(), f), func(v tines.AuditLog) int { return v.Id })
		}},
		{"teams", "/api/v1/teams", "teams", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListTeams(context.Background(), f), func(v tines.Team) int { return v.ID })
		}},
		{"members", "/api/v1/teams/7/members", "members", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListTeamMembers(context.Background(), 7, f), func(v tines.TeamMember) int { return v.ID })
		}},
		{"users", "/api/v1/admin/users", "admin/users", func(c *tines.Client, f tines.ListFilter) iter.Seq2[int, error] {
			return listIDs(c.ListUsers(context.Background(), f), func(v tines.User) int { return v.Id })
		}},
	}
	cases := []struct {
		name         string
		limit        int
		stopAfter    int
		failurePage  int
		invalidJSON  bool
		wantIDs      []int
		wantRequests int
	}{
		{name: "stop_after_first_item", stopAfter: 1, wantIDs: []int{1}, wantRequests: 1},
		{name: "stop_at_page_boundary", stopAfter: 2, wantIDs: []int{1, 2}, wantRequests: 1},
		{name: "later_http_error", failurePage: 2, wantIDs: []int{1, 2}, wantRequests: 2},
		{name: "later_decode_error", failurePage: 2, invalidJSON: true, wantIDs: []int{1, 2}, wantRequests: 2},
		{name: "first_page_error", failurePage: 1, wantRequests: 1},
		{name: "limit_within_page", limit: 1, wantIDs: []int{1}, wantRequests: 1},
		{name: "limit_at_page_boundary", limit: 2, wantIDs: []int{1, 2}, wantRequests: 1},
		{name: "limit_on_second_page", limit: 3, wantIDs: []int{1, 2, 3}, wantRequests: 2},
		{name: "unlimited", wantIDs: []int{1, 2, 3, 4}, wantRequests: 2},
		{name: "limit_exceeds_results", limit: 10, wantIDs: []int{1, 2, 3, 4}, wantRequests: 2},
	}
	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					requests := 0
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						requests++
						assert.Equal(t, endpoint.path, r.URL.Path)
						page := 1
						if r.URL.Query().Get("page") == "2" {
							page = 2
						}
						if page == tc.failurePage {
							if tc.invalidJSON {
								_, _ = w.Write([]byte("invalid json"))
							} else {
								w.WriteHeader(http.StatusServiceUnavailable)
							}
							return
						}
						next := `"next_page":null,"next_page_number":null`
						if page == 1 {
							next = fmt.Sprintf(`"next_page":"http://%s%s?page=2&per_page=2","next_page_number":2`, r.Host, endpoint.path)
						}
						_, _ = fmt.Fprintf(w, `{"%s":[{"id":%d},{"id":%d}],"meta":{%s,"per_page":2,"pages":2,"count":4}}`, endpoint.key, page*2-1, page*2, next)
					}))
					defer ts.Close()
					cli, err := tines.NewClient(tines.SetApiKey("foo"), tines.SetTenantUrl(ts.URL))
					require.NoError(t, err)
					var ids []int
					errors := 0
					for id, err := range endpoint.list(cli, tines.NewListFilter(tines.WithMaxResults(tc.limit))) {
						if err != nil {
							errors++
							continue
						}
						ids = append(ids, id)
						if len(ids) == tc.stopAfter {
							break
						}
					}
					assert.Equal(t, tc.wantIDs, ids)
					assert.Equal(t, tc.wantRequests, requests)
					if tc.failurePage > 0 {
						assert.Equal(t, 1, errors)
					} else {
						assert.Zero(t, errors)
					}
				})
			}
		})
	}
}
