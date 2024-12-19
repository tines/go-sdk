package tines

import (
	"context"
	"encoding/json"
	"net/http"
)

type TenantInfo struct {
	Stack StackInfo `json:"stack,omitempty"`
}

type StackInfo struct {
	Name      string   `json:"name,omitempty"`
	Type      string   `json:"type,omitempty"`
	Region    string   `json:"region,omitempty"`
	EgressIps []string `json:"egress_ips,omitempty"`
}

type WorkerStats struct {
	CurrentWorkers int `json:"current_workers,omitempty"`
	MaxWorkers     int `json:"max_workers,omitempty"`
	QueueCount     int `json:"queue_count,omitempty"`
	QueueLatency   int `json:"queue_latency,omitempty"`
}

func (c *Client) GetInfo(ctx context.Context) (*TenantInfo, error) {
	t := TenantInfo{}
	resource := "/api/v1/info"

	body, err := c.doRequest(ctx, http.MethodGet, resource, nil, nil)
	if err != nil {
		return &t, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return &t, err
	}

	return &t, nil
}

func (c *Client) GetWorkerStats(ctx context.Context) (*WorkerStats, error) {
	w := WorkerStats{}
	resource := "/api/v1/info/worker_stats"

	body, err := c.doRequest(ctx, http.MethodGet, resource, nil, nil)
	if err != nil {
		return &w, err
	}

	err = json.Unmarshal(body, &w)
	if err != nil {
		return &w, err
	}

	return &w, nil
}
