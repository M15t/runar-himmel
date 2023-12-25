package contextutil

import (
	"github.com/labstack/echo/v4"
)

// MWContext is a middleware that injects custom context into echo.Context
func MWContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &HTTPContext{Context: c}
			cc.SetAuthUser()

			return next(cc)
		}
	}
}
