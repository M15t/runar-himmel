package contextutil

import (
	"context"
	"runar-himmel/internal/types"
)

// AuthUser returns authenticated user info
func (h *HTTPContext) AuthUser() *types.AuthUser {
	return h.au
}

// RealIP returns real ip address
func (h *HTTPContext) RealIP() string {
	return h.Context.RealIP()
}

// UserAgent returns user agent
func (h *HTTPContext) UserAgent() string {
	return h.Context.Request().UserAgent()
}

// GetContext returns context
func (h *HTTPContext) GetContext() context.Context {
	return h.Context.Request().Context()
}

// SetAuthUser sets authenticated user info
func (h *HTTPContext) SetAuthUser() {
	// * add more if needed
	// * note: that need to modidy internal/types/auth.go AuthUser struct as well

	h.au = &types.AuthUser{
		ID:    h.getValue("id"),
		Name:  h.getValue("name"),
		Email: h.getValue("email"),
		Role:  h.getValue("role"),
		// Add more fields if needed
	}
}

// getValue safely retrieves a string value from the context
func (h *HTTPContext) getValue(key string) string {
	if value, ok := h.Context.Get(key).(string); ok {
		return value
	}
	return ""
}
