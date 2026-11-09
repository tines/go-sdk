package tines_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tines/go-sdk/tines"
)

func TestResourceAndCredentialHTTPErrors(t *testing.T) {
	methods := []struct {
		name   string
		method string
		path   string
		call   func(*tines.Client) error
	}{
		{"GetResource", http.MethodGet, "/api/v1/global_resources/1", func(c *tines.Client) error {
			_, err := c.GetResource(context.Background(), 1)
			return err
		}},
		{"UpdateResource", http.MethodPut, "/api/v1/global_resources/1", func(c *tines.Client) error {
			_, err := c.UpdateResource(context.Background(), 1, &tines.Resource{Id: 1, Value: "updated"})
			return err
		}},
		{"AppendResourceElement", http.MethodPost, "/api/v1/global_resources/1/append", func(c *tines.Client) error {
			_, err := c.AppendResourceElement(context.Background(), 1, &tines.ResourceElement{ResourceId: 1, Value: "appended"})
			return err
		}},
		{"RemoveResourceElement", http.MethodPost, "/api/v1/global_resources/1/remove", func(c *tines.Client) error {
			_, err := c.RemoveResourceElement(context.Background(), 1, &tines.ResourceElement{ResourceId: 1, Key: "example"})
			return err
		}},
		{"ReplaceResourceElement", http.MethodPost, "/api/v1/global_resources/1/replace", func(c *tines.Client) error {
			_, err := c.ReplaceResourceElement(context.Background(), 1, &tines.ResourceElement{ResourceId: 1, Key: "example", Value: "replacement"})
			return err
		}},
		{"GetCredential", http.MethodGet, "/api/v1/user_credentials/1", func(c *tines.Client) error {
			_, err := c.GetCredential(context.Background(), 1)
			return err
		}},
		{"UpdateCredential", http.MethodPut, "/api/v1/user_credentials/1", func(c *tines.Client) error {
			_, err := c.UpdateCredential(context.Background(), 1, &tines.Credential{Id: 1, Mode: tines.CredentialTypeText})
			return err
		}},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, status := range []int{http.StatusNotFound, http.StatusForbidden, http.StatusUnprocessableEntity, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable} {
				t.Run(fmt.Sprint(status), func(t *testing.T) {
					requests := 0
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						requests++
						assert.Equal(t, method.method, r.Method)
						assert.Equal(t, method.path, r.URL.Path)
						w.WriteHeader(status)
						_, _ = w.Write([]byte(`{"errors":[{"message":"API failure","details":"original details"}]}`))
					}))
					defer ts.Close()
					cli, err := tines.NewClient(tines.SetApiKey("foo"), tines.SetTenantUrl(ts.URL))
					require.NoError(t, err)
					err = method.call(cli)
					var apiErr tines.Error
					require.ErrorAs(t, err, &apiErr)
					assert.Equal(t, 1, requests)
					assert.Equal(t, status, apiErr.StatusCode)
					wantType := tines.ErrorTypeRequest
					if status >= http.StatusInternalServerError {
						wantType = tines.ErrorTypeServer
					}
					assert.Equal(t, wantType, apiErr.Type)
					assert.Equal(t, []tines.ErrorMessage{{Message: "API failure", Details: "original details"}}, apiErr.Errors)
				})
			}
		})
	}
}
