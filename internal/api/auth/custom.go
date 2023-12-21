package auth

import (
	"net/http"
	"runar-himmel/pkg/server"
)

// Custom errors
var (
	ErrInvalidCredentials  = server.NewHTTPError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email or password is incorrect")
	ErrUserBlocked         = server.NewHTTPError(http.StatusUnauthorized, "USER_BLOCKED", "Your account has been blocked and may not login")
	ErrInvalidRefreshToken = server.NewHTTPError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
	ErrTokenExpired        = server.NewHTTPError(http.StatusUnauthorized, "TOKEN_EXPIRED", "Invalid refresh token")
	ErrInvalidPayloadType  = server.NewHTTPError(http.StatusUnauthorized, "INVALID_PAYLOAD_TYPE", "Invalid payload type")
	ErrRefreshToken        = server.NewHTTPError(http.StatusInternalServerError, "REFRESH_TOKEN_ERROR", "An error occur while refreshing token")
	ErrInvalidGrantType    = server.NewHTTPError(http.StatusBadRequest, "INVALID_GRANT_TYPE", "Invalid grant type")
)
