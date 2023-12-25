package requestutil

// ListCondition represents a generic type for listing and filtering data
type ListCondition[T any] struct {
	// Current page number, starts from 1
	Page int
	// Number of records per page
	PerPage int
	// Field names for sorting, separated by comma.
	// The direction is specified by `+` (ASC) and `-` (DESC) prefix, eg: `+name,-age`
	// WARNING: SQL Injection vulnerability! User input must be validated before sending to database
	Sort string
	// Whether to count the total records. If not, the returning `count` number will always be zero.
	Count bool
	// Custom filter type
	Filter T
}

// ListQueryRequest contains basic information to do sorting and paging for a listing api
// WARNING: SQL Injection vulnerability! Sort param must be validated before sending to database
// swagger:model
type ListQueryRequest struct {
	// Current page number, starts from 1
	// in: query
	// default: 1
	Page int `json:"p" query:"p"`
	// Number of records per page
	// in: query
	// default: 25
	PerPage int `json:"pp" query:"pp"`
	// Field names for sorting, separated by comma.
	// The direction is specified by `+` (ASC) and `-` (DESC) prefix, eg: `+name,-age`
	// in: query
	Sort string `json:"s" query:"s"`
}

// ListQueryCondition represents a generic type for listing and filtering data
type ListQueryCondition ListCondition[[]any]
