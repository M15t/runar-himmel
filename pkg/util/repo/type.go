package repoutil

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

// ListQueryCondition represents a generic type for listing and filtering data
type ListQueryCondition ListCondition[[]any]
