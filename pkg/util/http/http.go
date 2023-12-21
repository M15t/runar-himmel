package httputil

import (
	"runar-himmel/pkg/server"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

// ReqID returns id url parameter.
func ReqID(c echo.Context, param ...string) (id string, err error) {
	name := "id"
	if len(param) > 0 {
		name = param[0]
	}
	id = c.Param(name)
	_, err = ulid.ParseStrict(id)
	if err != nil {
		return "", server.NewHTTPValidationError("Record not found")
	}
	return
}

// ReqIDint64 returns id url parameter. Support int64
func ReqIDint64(c echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, server.NewHTTPValidationError("Invalid ID")
	}
	return id, nil
}

// ReqIDint returns id url parameter.
func ReqIDint(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, server.NewHTTPValidationError("Invalid ID")
	}
	return id, nil
}

// TrimSpacePointer trims leading and trailing spaces from a pointer string
func TrimSpacePointer(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

// RemoveSpacePointer remove all spaces from a pointer string
func RemoveSpacePointer(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.Replace(*s, " ", "", -1)
	return &trimmed
}

// // ListRequest holds data of listing request from react-admin
// // Note: To add these parameters to swagger:operation, check the file cmd/api/swagger.go
// // swagger:ignore
// type ListRequest struct {
// 	// Number of records per page
// 	// default: 25
// 	Limit int `json:"l,omitempty" query:"l"`
// 	// Current page number
// 	// default: 1
// 	Page int `json:"p,omitempty" query:"p"`
// 	// Field name for sorting
// 	// default:
// 	Sort string `json:"s,omitempty" query:"s"`
// 	// Sort direction, must be one of ASC, DESC
// 	// default:
// 	Order string `json:"o,omitempty" query:"o"`
// 	// JSON string of filter. E.g: {"field_name":"value"}
// 	// default:
// 	Filter string `json:"f,omitempty" query:"f"`
// }

// // ReqListQuery parses url query string for listing request
// func ReqListQuery(c echo.Context) (*dbutil.ListQueryCondition, error) {
// 	isValidParams := regexp.MustCompile(`^[a-zA-Z0-9._"]*$`).MatchString

// 	lr := &ListRequest{}
// 	if err := c.Bind(lr); err != nil {
// 		return nil, err
// 	}

// 	lq := &dbutil.ListQueryCondition{
// 		Page:    lr.Page,
// 		PerPage: lr.Limit,
// 		Filter:  gowhere.WithConfig(gowhere.Config{Strict: true}),
// 	}

// 	if lr.Filter != "" {
// 		var filter interface{}
// 		err := json.Unmarshal([]byte(lr.Filter), &filter)
// 		if err != nil {
// 			return nil, server.NewHTTPValidationError("Invalid filter, expecting JSON string").SetInternal(err)
// 		}

// 		if err := lq.Filter.Where(filter).Build().Error; err != nil {
// 			return nil, server.NewHTTPValidationError("Cannot parse filter").SetInternal(err)
// 		}
// 	}

// 	if lr.Sort != "" {
// 		if !isValidParams(lr.Sort) || len(lr.Sort) > 50 {
// 			return nil, server.NewHTTPValidationError("Invalid params for sort")
// 		}

// 		sortField := lr.Sort
// 		sortOrder := "ASC" // default
// 		if lr.Order != "" && strings.ToLower(lr.Order) == "desc" {
// 			sortOrder = "DESC"
// 		}
// 		lq.Sort = []string{sortField + " " + sortOrder}
// 	}

// 	return lq, nil
// }
