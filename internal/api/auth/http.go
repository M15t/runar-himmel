package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"runar-himmel/internal/types"
	"runar-himmel/pkg/server"
)

// HTTP represents auth http service
type HTTP struct {
	svc Service
}

// Service represents auth service interface
type Service interface {
	Login(echo.Context, Credentials) (*types.AuthToken, error)
	RefreshToken(echo.Context, RefreshTokenData) (*types.AuthToken, error)
}

// NewHTTP attaches handlers to Echo routers under given group
func NewHTTP(svc Service, eg *echo.Group) {
	h := HTTP{svc: svc}

	// swagger:operation POST /auth/login auth authLogin
	// ---
	// summary: Logs in user by email, password and grant_type
	// security: []
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body. `grant_type` should be `app` or `portal`
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Credentials"
	// responses:
	//   "200":
	//     description: Access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   default:
	//     description: 'Possible errors: 400, 401, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.POST("/login", h.login)

	// swagger:operation POST /auth/refresh-token auth authRefreshToken
	// ---
	// summary: Refresh access token
	// security: []
	// parameters:
	// - name: token
	//   in: body
	//   description: The given `refresh_token` when login
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RefreshTokenData"
	// responses:
	//   "200":
	//     description: New access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   default:
	//     description: 'Possible errors: 400, 401,500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.POST("/refresh-token", h.refreshToken)
}

func (h *HTTP) login(c echo.Context) error {
	r := Credentials{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	if r.Email == "" {
		r.Email = strings.ToLower(strings.TrimSpace(r.Username))
	}

	if !lo.Contains([]string{
		"app", "portal",
	}, r.GrantType) {
		return server.NewHTTPValidationError("Invalid context")
	}

	resp, err := h.svc.Login(c, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) refreshToken(c echo.Context) error {
	r := RefreshTokenData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	resp, err := h.svc.RefreshToken(c, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
