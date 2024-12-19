package paginate

import "net/url"

type Meta struct {
	CurrentPage  string `json:"current_page,omitempty"`
	PreviousPage string `json:"previous_page,omitempty"`
	NextPage     string `json:"next_page"`
	NextPageNum  int    `json:"next_page_number"`
	PerPage      int    `json:"per_page,omitempty"`
	Pages        int    `json:"pages,omitempty"`
	Count        int    `json:"count,omitempty"`
}

type Cursor struct {
	Meta           Meta
	totalReturned  int
	TotalRequested int
}

func (c *Cursor) MoreResultsAvailable() bool {
	return c.Meta.NextPageNum > 0
}

func (c *Cursor) IncrementCounter() {
	c.totalReturned++
}

func (c *Cursor) CurrentCounter() int {
	return c.totalReturned
}

func (c *Cursor) ReturnMoreResults() bool {
	return c.MoreResultsAvailable() && !c.MaxResultsReturned()
}

func (c *Cursor) MaxResultsReturned() bool {
	if c.TotalRequested == 0 {
		return false
	}
	return c.totalReturned >= c.TotalRequested
}

func (c *Cursor) UpdatePagination(m Meta) {
	c.Meta = m
}

func (c *Cursor) GetNextPageParams() map[string]any {
	uri, _ := url.Parse(c.Meta.NextPage)
	params := make(map[string]any)
	for k, v := range uri.Query() {
		params[k] = v[0]
	}
	return params
}
