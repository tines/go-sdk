package tines_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

func TestClientSuccess(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil)
	defer ts.Close()

	// Validate that the Tines CLI gets instantiated correctly.
	_, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

}

func TestClientError(t *testing.T) {
	assert := assert.New(t)

	// Validate that we throw an error when required params are missing.
	_, err := tines.NewClient()
	assert.Error(err)
}

func createTestServer(assert *assert.Assertions, respStatus int, respBody []byte) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate that the client is sending expected request values
		assert.Equal("application/json", r.Header.Get("Content-Type"), "client should send JSON data")
		assert.Equal("application/json", r.Header.Get("Accept"))
		assert.Equal("Bearer foo", r.Header.Get("Authorization"))
		assert.Equal("Tines/GoSdk", r.Header.Get("User-Agent"))

		w.WriteHeader(respStatus)

		w.Write(respBody) //nolint:errcheck
	}))
	return ts
}
