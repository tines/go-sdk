package tines

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type AuditLog struct {
	CreatedAt     string      `json:"created_at"`
	OperationName string      `json:"operation_name"`
	Id            int         `json:"id"`
	Inputs        interface{} `json:"inputs"`
	Outputs       interface{} `json:"outputs"`
	RequestIP     string      `json:"request_ip"`
	RequestUA     string      `json:"request_user_agent"`
	StoryID       int         `json:"story_id"`
	TenantID      int         `json:"tenant_id"`
	UpdatedAt     string      `json:"updated_at"`
	UserEmail     string      `json:"user_email"`
	UserID        int         `json:"user_id"`
	UserName      string      `json:"user_name"`
}

type AuditLogList struct {
	AuditLogs []AuditLog    `json:"audit_logs"`
	Meta      paginate.Meta `json:"meta"`
}

func (c *Client) ListAuditLogs(ctx context.Context, f *ListFilter) iter.Seq2[AuditLog, error] {
	var auditlogList, resultList AuditLogList
	resource := "/api/v1/audit_logs"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}

	return func(yield func(AuditLog, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(AuditLog{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(AuditLog{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.AuditLogs {
				auditlogList.AuditLogs = append(auditlogList.AuditLogs, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = AuditLogList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range auditlogList.AuditLogs {
			if !yield(v, nil) {
				return
			}
		}
	}
}
