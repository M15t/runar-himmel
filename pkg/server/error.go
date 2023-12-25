package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const (
	// InternalErrorType type of common errors
	InternalErrorType = "INTERNAL"
	// GenericErrorType type of common errors
	GenericErrorType = "GENERIC"
	// ValidationErrorType type of common errors
	ValidationErrorType = "VALIDATION"
)

// ErrorResponse represents the error response
// swagger:model
type ErrorResponse struct {
	Error *HTTPError `json:"error"`
}

// HTTPError represents an error that occurred while handling a request
type HTTPError struct {
	Code     int    `json:"code"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}

// NewHTTPError creates a new HTTPError instance
func NewHTTPError(code int, etype string, message ...string) *HTTPError {
	he := &HTTPError{Code: code, Type: etype}
	if len(message) > 0 {
		he.Message = message[0]
	} else {
		he.Message = http.StatusText(code)
	}
	return he
}

// NewHTTPInternalError creates a new HTTPError instance for internal error
func NewHTTPInternalError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusInternalServerError, Type: InternalErrorType, Message: message}
}

// NewHTTPGenericError creates a new HTTPError instance for generic error
func NewHTTPGenericError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Type: GenericErrorType, Message: message}
}

// NewHTTPValidationError creates a new HTTPError instance for validation error
func NewHTTPValidationError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Type: ValidationErrorType, Message: message}
}

// Error makes it compatible with `error` interface
func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, type=%s, message=%s", he.Code, he.Type, he.Message)
	}
	return fmt.Sprintf("code=%d, type=%s, message=%s, internal=%v", he.Code, he.Type, he.Message, he.Internal)
}

// SetInternal sets actual internal error for more details
func (he *HTTPError) SetInternal(err error) *HTTPError {
	he.Internal = err
	return he
}

// MarshalJSON modification
func (he *HTTPError) MarshalJSON() ([]byte, error) {
	type alias HTTPError
	type custom struct {
		*alias
		Internal string `json:"internal,omitempty"`
	}
	output := custom{alias: (*alias)(he)}
	if he.Internal != nil {
		output.Internal = he.Internal.Error()
	}
	return json.Marshal(output)
}

// ErrorHandler represents the custom http error handler
type ErrorHandler struct {
	e *echo.Echo
}

// NewErrorHandler returns the ErrorHandler instance
func NewErrorHandler(e *echo.Echo) *ErrorHandler {
	return &ErrorHandler{e}
}

// Handle is a centralized HTTP error handler.
func (ce *ErrorHandler) Handle(err error, c echo.Context) {
	httpErr := NewHTTPError(http.StatusInternalServerError, InternalErrorType)

	switch e := err.(type) {
	case *HTTPError:
		if e.Code != 0 {
			httpErr.Code = e.Code
		}
		if e.Type != "" {
			httpErr.Type = e.Type
		} else {
			httpErr.Type = GenericErrorType
		}
		if e.Message != "" {
			httpErr.Message = e.Message
		}
		if e.Internal != nil {
			ce.e.Logger.Errorf("internal err: %+v", e.Internal)
			// httpErr.Internal = e.Internal
		}

	case *echo.HTTPError:
		httpErr.Code = e.Code
		httpErr.Type = GenericErrorType
		switch em := e.Message.(type) {
		case string:
			httpErr.Message = em
		case []string:
			httpErr.Message = strings.Join(em, "\n")
		case map[string]interface{}:
			if jsonStr, err := json.Marshal(em); err == nil {
				httpErr.Message = string(jsonStr)
			}
		default:
			httpErr.Message = fmt.Sprintf("%+v", em)
		}
		if e.Internal != nil {
			ce.e.Logger.Errorf("internal err: %+v", e.Internal)
			// httpErr.Internal = e.Internal
		}

	case validator.ValidationErrors:
		httpErr.Code = http.StatusBadRequest
		httpErr.Type = ValidationErrorType
		var errMsg []string
		for _, v := range e {
			errMsg = append(errMsg, getVldErrorMsg(v))
		}
		httpErr.Message = strings.Join(errMsg, "\n")
	default:
		if ce.e.Debug {
			httpErr.Message = err.Error()
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(httpErr.Code)
		} else {
			if !ce.e.Debug && httpErr.Code == http.StatusInternalServerError {
				ce.e.Logger.Error(httpErr)
				httpErr.Message = http.StatusText(httpErr.Code)
				httpErr.Internal = nil
				httpErr.Type = GenericErrorType
			}
			err = c.JSON(httpErr.Code, ErrorResponse{Error: httpErr})
		}
		if err != nil {
			ce.e.Logger.Error(err)
		}
	}
}

var validationErrors = map[string]string{
	"required":   " is required, but was not received",
	"min":        "'s value or length is less than allowed",
	"max":        "'s value or length is bigger than allowed",
	"date":       " should be a valid date in form of YYYY-MM-DD",
	"email":      " is an invalid email address",
	"phone":      " is an invalid phone number",
	"url":        " is an invalid URL",
	"coordinate": " is an invalid latitude, longitude coordinate",
}

func getVldErrorMsg(v validator.FieldError) string {
	field := v.Field()
	vtag := v.ActualTag()
	vtagVal := v.Param()

	if msg, ok := validationErrors[vtag]; ok {
		return field + msg
	}

	switch vtag {
	case "min":
		return field + " minimum allowed length is " + vtagVal
	case "max":
		return field + " exceeds the maximum allowed length (" + vtagVal + ")"
	case "lte":
		return field + "'s length should be less than or equal to " + vtagVal
	case "lt":
		return field + "'s length should be less than " + vtagVal
	case "gte":
		return field + "'s length should be greater than or equal to " + vtagVal
	case "gt":
		return field + "'s length should be greater than " + vtagVal
	case "oneof":
		return field + " should be one of " + strings.Replace(vtagVal, " ", ", ", -1)
	case "ltfield":
		return field + " should be less than " + vtagVal
	case "gtfield":
		return field + " should be greater than " + vtagVal
	case "eqfield":
		return field + " does not match " + vtagVal
	case "datetime":
		return field + "'s format should be: " + vtagVal
	}

	return field + " failed on " + vtag + " validation"
}
