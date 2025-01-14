package tines_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

// API response value as of 2025-01-14.
const testAuditLogsResp = `
{
    "audit_logs": [
        {
            "created_at": "2025-01-11T03:57:28Z",
            "operation_name": "Action",
            "id": 1,
            "inputs": {
                "inputs": {
                    "actionId": 1
                }
            },
            "outputs": {},
            "request_ip": "1.1.1.1",
            "request_user_agent": "Foo",
            "tenant_id": 1,
            "updated_at": "2025-01-11T03:57:28Z",
            "user_email": "user@example.com",
            "user_id": 1,
            "user_name": "Example User",
            "story_id": 1
        }
    ],
    "meta": {
        "current_page": "https://example.tines.com/api/v1/audit_logs?per_page=1&page=1",
        "previous_page": null,
        "next_page": "https://example.tines.com/api/v1/audit_logs?per_page=1&page=2",
        "next_page_number": 2,
        "per_page": 1,
        "pages": 1,
        "count": 2
    }
}`

func TestAuditLogsList(t *testing.T) {
	assert := assert.New(t)

	ts := createTestServer(assert, http.StatusOK, []byte(testAuditLogsResp))
	defer ts.Close()

	cli, err := tines.NewClient(
		tines.SetApiKey("foo"),
		tines.SetTenantUrl(ts.URL),
	)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")
	if err != nil {
		return
	}

	ctx := context.Background()

	for l, err := range cli.ListAuditLogs(ctx, tines.NewListFilter()) {
		assert.Equal("Action", l.OperationName)
		assert.Nil(err)
	}

}
