package paging

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

func MovieSortableAttributes() *SortableAttributes {
	return newSortableAttributes([]string{
		"title", "released", "imdbRating", "score",
	})
}

func PersonSortableAttributes() *SortableAttributes {
	return newSortableAttributes([]string{
		"name", "born", "movieCount",
	})
}

func RatingSortableAttributes() *SortableAttributes {
	return newSortableAttributes([]string{
		"rating", "timestamp",
	})
}

type SortableAttributes struct {
	defaultValue string
	values       []string
}

func newSortableAttributes(values []string) *SortableAttributes {
	defaultValue := values[0]
	sort.Strings(values)
	return &SortableAttributes{defaultValue: defaultValue, values: values}
}

func (sa *SortableAttributes) contains(s string) bool {
	i := sort.SearchStrings(sa.values, s)
	return i < len(sa.values) && sa.values[i] == s
}

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

func ParsePaging(req *http.Request, sortableAttributes *SortableAttributes) *Paging {
	query := req.URL.Query()
	sortParameter := query.Get("sort")
	if !sortableAttributes.contains(sortParameter) {
		sortParameter = sortableAttributes.defaultValue
	}
	return &Paging{
		query: query.Get("q"),
		sort:  sortParameter,
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

func NewPaging(query string, sort string, order string, skip int, limit int) *Paging {
	return &Paging{
		query: query,
		sort:  sort,
		order: order,
		skip:  skip,
		limit: limit,
	}

}
