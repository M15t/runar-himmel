package contextutil

import (
	"context"
	"runar-himmel/internal/types"

	"github.com/labstack/echo/v4"
)

// HTTPContext is a custom context
type HTTPContext struct {
	echo.Context
	// User contains authenticated user info
	au *types.AuthUser
	// * add more if needed
}

// Context is a custom context interface
type Context interface {
	GetContext() context.Context
	AuthUser() *types.AuthUser
	RealIP() string
	UserAgent() string
}

// ensure it implements the Context interface
var _ Context = &HTTPContext{}

// NewContext returns new custom context
func NewContext(c echo.Context) Context {
	return c.(Context)
}
