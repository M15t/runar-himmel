package repoutil

import (
	"strings"

	"github.com/imdatngo/gowhere"
	"gorm.io/gorm"
)

func (r *Repo[T]) quoteCol(name string) string {
	b := &strings.Builder{}
	r.GDB.QuoteTo(b, name)
	return b.String()
}

// parseConds returns standard [sqlString, vars] format for query, powered by gowhere package (with default config)
func parseConds(conds []any) []any {
	if len(conds) == 1 {
		switch c := conds[0].(type) {
		case string:
			if len(c) == 0 {
				return []any{}
			}
			// Using parameterized query for string condition to prevent SQL injection
			return []any{"id = ?", c}
		case map[string]any:
			if len(c) == 0 {
				return []any{}
			}
			// Check if the map contains non-simple conditions (e.g., gowhere operators)
			if hasNonSimpleConditions(c) {
				// Use gowhere for complex conditions
				plan := gowhere.Where(c)
				return append([]any{plan.SQL()}, plan.Vars()...)
			}
			// Use simple equality conditions for Gorm
			return mapToGormConditions(c)
		case []any:
			if len(c) == 0 {
				return []any{}
			}
			// Use gowhere for slice conditions
			plan := gowhere.Where(c)
			return append([]any{plan.SQL()}, plan.Vars()...)
		case *gowhere.Plan:
			// Use gowhere plan directly
			return append([]any{c.SQL()}, c.Vars()...)
		case nil:
			return []any{}
		}
	}
	return conds
}

// parseSortValue returns the column name and direction for sorting value like +column or -column
func parseSortValue(s string) (col, dir string) {
	col = strings.TrimSpace(s)
	if col == "" || col == "+" || col == "-" {
		return "", ""
	}

	sign, col := col[0], col[1:]
	switch sign {
	case '+':
		dir = "ASC"
	case '-':
		dir = "DESC"
	default:
		col = string(sign) + col
		dir = "ASC"
	}
	return
}

// parseSortParam returns list of [column, direction] from comma separated sorting param
func parseSortParam(s string) [][]string {
	values := strings.Split(s, ",")
	l := make([][]string, 0, len(values))
	for _, v := range values {
		col, dir := parseSortValue(v)
		if col != "" {
			l = append(l, []string{col, dir})
		}
	}
	return l
}

// An util function to set pagination conditions
func withPaging(db *gorm.DB, page, perPage int) *gorm.DB {
	if perPage > 0 {
		db = db.Limit(perPage)
		if page > 1 {
			db = db.Offset(page*perPage - perPage)
		}
	}
	return db
}

// An util function to set sorting conditions
// WARNING: SQL Injection vulnerability! `quoteCol` function must take care of quoting column name properly
func withSorting(db *gorm.DB, sort string, quoteCol func(name string) string) *gorm.DB {
	if sort != "" {
		sp := parseSortParam(sort)
		values := []string{}
		for _, v := range sp {
			col := quoteCol(v[0])
			dir := v[1]
			values = append(values, col+" "+dir)
		}
		db = db.Order(strings.Join(values, ", "))
	}
	return db
}

// hasNonSimpleConditions checks if a map contains non-simple conditions (e.g., gowhere operators).
func hasNonSimpleConditions(m map[string]any) bool {
	for _, v := range m {
		switch v := v.(type) {
		case map[string]any:
			// Check if the value is a nested map
			if hasNonSimpleConditions(v) {
				return true
			}
		default:
			// Check if the value can be processed by gowhere.Where
			if gowhere.Where(v) != nil {
				return true
			}
		}
	}
	return false
}

// mapToGormConditions converts a map with simple equality conditions to Gorm conditions.
func mapToGormConditions(m map[string]any) []any {
	conds := make([]any, 0, len(m)*2)
	for k, v := range m {
		conds = append(conds, k+" = ?", v)
	}
	return conds
}
