package tines_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

func TestClientSuccess(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(assert, http.StatusOK, nil, nil)
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

func createTestServer(assert *assert.Assertions, expectRespStatus int, expectReqBody, expectRespBody []byte) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate that the client is sending expected request values
		assert.Equal("application/json", r.Header.Get("Content-Type"), "client should send JSON data")
		assert.Equal("application/json", r.Header.Get("Accept"))
		assert.Equal("Bearer foo", r.Header.Get("Authorization"))
		assert.Equal("Tines/GoSdk", r.Header.Get("User-Agent"))

		// Optionally validate that the request body from the client matches a particular format
		if expectReqBody != nil {
			defer r.Body.Close()

			reqBody, err := io.ReadAll(r.Body)
			assert.Nil(err, "HTTP request body should be readable")

			assert.JSONEq(string(expectReqBody), string(reqBody), "HTTP request body should match expected payload")
		}

		w.WriteHeader(expectRespStatus)

		w.Write(expectRespBody) //nolint:errcheck
	}))
	return ts
}
