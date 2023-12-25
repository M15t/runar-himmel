package requestutil

// ToListQueryCond transforms the service request to list query conditions
func (l *ListQueryRequest) ToListQueryCond(filter []any) *ListQueryCondition {
	return &ListQueryCondition{
		Page:    l.Page,
		PerPage: l.PerPage,
		Sort:    l.Sort,
		Count:   true,
		Filter:  filter,
	}
}
