package paginate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/internal/paginate"
)

var TestCursorWithLimit = paginate.Cursor{
	Meta: paginate.Meta{
		CurrentPage:  "https://example.com/?per_page=1&page=1",
		PreviousPage: "",
		NextPage:     "https://example.com?per_page=1&page=2",
		NextPageNum:  2,
		PerPage:      1,
		Pages:        3,
		Count:        3,
	},
	TotalRequested: 2,
}

var TestCursorNoLimit = paginate.Cursor{
	Meta: paginate.Meta{
		CurrentPage:  "https://example.com/?per_page=1&page=1",
		PreviousPage: "",
		NextPage:     "https://example.com?per_page=1&page=2",
		NextPageNum:  2,
		PerPage:      1,
		Pages:        2,
		Count:        2,
	},
	TotalRequested: 0,
}

func TestPaginationWithLimit(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, TestCursorWithLimit.MoreResultsAvailable(), "indicate that more results are available if next page is not nil")
	assert.Equal(false, TestCursorWithLimit.MaxResultsReturned(), "indicate that the maximum requested number of results has not been returned yet")
	assert.Equal(true, TestCursorWithLimit.ReturnMoreResults(), "indicate that more results should be returned")
	assert.Equal(0, TestCursorWithLimit.CurrentCounter(), "indicate that no results have been returned yet")

	TestCursorWithLimit.IncrementCounter()

	assert.Equal(1, TestCursorWithLimit.CurrentCounter(), "counter should increment by one")

	params := TestCursorWithLimit.GetNextPageParams()

	assert.Equal(map[string]interface{}{"page": "2", "per_page": "1"}, params, "next page params should get successfully extracted")

	TestCursorWithLimit.IncrementCounter()

	assert.Equal(false, TestCursorWithLimit.ReturnMoreResults(), "indicate that the maximum number of requested results have been returned")
}

func TestPaginationNoLimit(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(false, TestCursorNoLimit.MaxResultsReturned(), "indicate that there is no maximum number of results to return")

	newMeta := paginate.Meta{
		// Next Page info will be nil in the server response, so we don't set it here.
		CurrentPage:  "https://example.com/?per_page=1&page=2",
		PreviousPage: "https://example.com?per_page=1&page=1",
		PerPage:      1,
		Pages:        2,
		Count:        2,
	}

	TestCursorNoLimit.UpdatePagination(newMeta)

	assert.Equal(false, TestCursorNoLimit.MoreResultsAvailable(), "indicate that no more results are available")
}
