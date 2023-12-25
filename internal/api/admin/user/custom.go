package user

import (
	"net/http"
	"runar-himmel/pkg/server"
)

// Custom errors
var (
	ErrIncorrectPassword = server.NewHTTPError(http.StatusBadRequest, "INCORRECT_PASSWORD", "Incorrect old password")
	ErrUserNotFound      = server.NewHTTPError(http.StatusBadRequest, "USER_NOTFOUND", "User not found")
	ErrEmailExisted      = server.NewHTTPValidationError("Email already existed")
)
