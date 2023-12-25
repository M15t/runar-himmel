package session

import (
	"net/http"
	"runar-himmel/pkg/server"
)

// Custom errors
var (
	ErrSessionNotFound = server.NewHTTPError(http.StatusBadRequest, "SESSION_NOTFOUND", "Session not found")
)
