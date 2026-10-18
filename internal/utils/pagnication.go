package utils

import "strconv"

// Pagination holds page info
type Pagination struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

// NewPagination parses pagination query params safely
func NewPagination(pageStr, limitStr, sort, order string) Pagination {
	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// Offset returns SQL OFFSET
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// OrderClause returns safe ORDER BY clause
func (p Pagination) OrderClause(allowedFields map[string]string, defaultField string) string {
	field := defaultField

	if p.Sort != "" {
		if mapped, ok := allowedFields[p.Sort]; ok {
			field = mapped
		}
	}

	return field + " " + p.Order
}