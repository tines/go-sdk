package tines_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/tines"
)

const (
	// Sanitized API response value as of 2025-01-14.
	testGetInfoResp = `
{
    "stack": {
        "name": "us1",
        "type": "shared",
        "region": "us-west-2",
        "egress_ips": [
            "1.1.1.1",
            "1.1.1.1"
        ]
    }
}`
	// Sanitized API response value as of 2025-01-14.
	testGetWorkerStatsResp = `
{
	"current_workers": 10,
	"max_workers": 20,
	"queue_count": 100,
	"queue_latency": 500
}`
)

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	ts := createTestServer(assert, 200, nil, []byte(testGetInfoResp))
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

	ti, err := cli.GetInfo(ctx)

	assert.Nil(err)
	assert.Equal("us-west-2", ti.Stack.Region)

}

func TestGetWorkerStats(t *testing.T) {
	assert := assert.New(t)

	ts := createTestServer(assert, 200, nil, []byte(testGetWorkerStatsResp))
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

	ws, err := cli.GetWorkerStats(ctx)

	assert.Nil(err)
	assert.Equal(10, ws.CurrentWorkers)

}
