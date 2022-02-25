package paging

import (
	"net/http"
	"net/url"
	"strconv"
)

type Paging struct {
	query string
	sort  string
	order string
	skip  int
	limit int
}

func (p Paging) Query() string {
	return p.query
}

func (p Paging) Sort() string {
	return p.sort
}

func (p Paging) Order() string {
	return p.order
}

func (p Paging) Skip() int {
	return p.skip
}

func (p Paging) Limit() int {
	return p.limit
}

func ParsePaging(req *http.Request) *Paging {
	query := req.URL.Query()
	return &Paging{
		query: query.Get("q"),
		// TODO: validate sort
		sort:  query.Get("sort"),
		order: query.Get("order"),
		skip:  getIntOrDefault(query, "skip", 0),
		limit: getIntOrDefault(query, "limit", 6),
	}
}

func getIntOrDefault(query url.Values, key string, defaultValue int) int {
	rawSkip := query.Get(key)
	if rawSkip == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(rawSkip)
	if err != nil {
		return defaultValue
	}
	return result
}
