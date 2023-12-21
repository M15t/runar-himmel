package secure

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSConfig represents secure specific CORS config
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
}

// Headers adds general security headers for basic security measures
func Headers(securityPolicy string) echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
		ContentSecurityPolicy: securityPolicy,
	})
}

// CORS adds Cross-Origin Resource Sharing support
func CORS(cfg CORSConfig) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  cfg.AllowOrigins,
		AllowMethods:  cfg.AllowMethods,
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Etag"},
		MaxAge:        86400,
	})
}

// BodyDump prints out the request body for debugging purpose
func BodyDump() echo.MiddlewareFunc {
	secretFields := []string{"password", "key", "token", "cert", "username", "email", "phone", "mobile"}
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			// request method
			requestMethod := c.Request().Method
			if requestMethod != http.MethodPatch && requestMethod != http.MethodPost && requestMethod != http.MethodPut {
				return true
			}
			// only support json content type
			reqContentType := c.Request().Header.Get("Content-Type")
			resContentType := c.Response().Header().Get("Content-Type")
			if !strings.Contains(reqContentType, "application/json") && !strings.Contains(resContentType, "application/json") {
				return true
			}
			// request too large
			if length := c.Request().ContentLength; length > 1000000 {
				c.Logger().Warnf("Skipped BodyDump, request body too large: %d", length)
				return true
			}
			return false
		},
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") && len(reqBody) > 0 {
				var bodymap map[string]interface{}
				if err := json.Unmarshal(reqBody, &bodymap); err == nil {
					bodymap = censorSecerts(c.Request().URL.String(), bodymap, secretFields)
					reqBody, _ = json.Marshal(bodymap)
					c.Logger().Info("Request Body: " + string(reqBody))
				}
			}

			if strings.Contains(c.Response().Header().Get("Content-Type"), "application/json") && len(resBody) > 0 {
				var bodymap map[string]interface{}
				if err := json.Unmarshal(resBody, &bodymap); err == nil {
					bodymap = censorSecerts(c.Request().URL.String(), bodymap, secretFields)
					resBody, _ = json.Marshal(bodymap)
					c.Logger().Info("Response Body: " + string(resBody))
				}
			}
		},
	})
}

// DisableCache sets the Cache-Control directive to no-store.
func DisableCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			c.Response().Header().Set("Cache-Control", "no-store")
			return next(c)
		}
	}
}

// SimpleCORS returns a CORS middleware with minimum configurations. Preflighted request is not allowed though.
func SimpleCORS(allowOrigins []string) echo.MiddlewareFunc {
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// Check allowed origins
			origin := req.Header.Get(echo.HeaderOrigin)
			allowed := ""
			for _, o := range allowOrigins {
				if o == "*" || o == origin {
					allowed = o
					break
				}
			}

			// Simple request
			switch req.Method {
			case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				res.Header().Add(echo.HeaderVary, echo.HeaderOrigin)
				res.Header().Set(echo.HeaderAccessControlAllowOrigin, allowed)
				return next(c)
			}

			// Preflight request is only allowed when "all" origins are allowed
			if req.Method == http.MethodOptions && allowed == "*" {
				res.Header().Add(echo.HeaderVary, echo.HeaderOrigin)
				res.Header().Add(echo.HeaderVary, echo.HeaderAccessControlRequestMethod)
				res.Header().Add(echo.HeaderVary, echo.HeaderAccessControlRequestHeaders)
				res.Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
				res.Header().Set(echo.HeaderAccessControlAllowMethods, "*")
				res.Header().Set(echo.HeaderAccessControlAllowHeaders, "*")
				return c.NoContent(http.StatusNoContent)
			}

			return echo.ErrMethodNotAllowed
		}
	}
}

func censorSecerts(uri string, body map[string]interface{}, secrets []string) map[string]interface{} {
	for key, val := range body {
		found := false
		lowerkey := strings.ToLower(key)
		for _, secretKey := range secrets {
			if secretKey == lowerkey || strings.Contains(lowerkey, secretKey) {
				found = true
				break
			}
		}
		if found {
			body[key] = "***"
			continue
		}

		switch v := val.(type) {
		case map[string]interface{}:
			body[key] = censorSecerts(uri, v, secrets)
			continue
		}
	}

	return body
}
