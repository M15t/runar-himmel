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

// AddFilter adds filter condition
func (lqc *ListQueryCondition) AddFilter(conds ...any) *ListQueryCondition {
	lqc.Filter = append(lqc.Filter, conds)
	return lqc
}

// parseConds returns standard [sqlString, vars] format for query, powered by gowhere package (with default config)
func parseConds(conds []any) []any {
	if len(conds) == 1 {
		var plan *gowhere.Plan
		switch c := conds[0].(type) {
		case string:
			if len(c) == 0 {
				return []any{}
			}
			// better safe than sorry
			return []any{"/* WARNING: SQL Injection vulnerability! */ id = ?", c}
		case map[string]any:
			if len(c) == 0 {
				return []any{}
			}
			//? this eliminates gorm's ability to query using map...
			plan = gowhere.Where(c)
		case []any:
			if len(c) == 0 {
				return []any{}
			}
			plan = gowhere.Where(c)
		case *gowhere.Plan:
			plan = c
		case nil:
			return []any{}
		}

		if plan != nil {
			return append([]any{plan.SQL()}, plan.Vars()...)
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

// A util function to set pagination conditions
func withPaging(db *gorm.DB, page, perPage int) *gorm.DB {
	if perPage > 0 {
		db = db.Limit(perPage)
		if page > 1 {
			db = db.Offset(page*perPage - perPage)
		}
	}
	return db
}

// A util function to set sorting conditions
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
