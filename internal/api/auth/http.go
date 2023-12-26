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

	// @Summary Logs in user by email, password and grant_type
	// @Description Request body. `grant_type` should be `app` or `portal`
	// @Accept  json
	// @Produce  json
	// @Param   request	 body    Credentials     true        "Request body"
	// @Success 200 {object} AuthToken
	// @Failure 400 {object} server.ErrorResponse
	// @Failure 401 {object} server.ErrorResponse
	// @Failure 500 {object} server.ErrorResponse
	// @Router /auth/login [post]
	eg.POST("/login", h.login)

	// @Summary Refresh access token
	// @Description Request body
	// @Accept  json
	// @Produce  json
	// @Param   request	 body    RefreshTokenData     true        "Request body"
	// @Success 200 {object} AuthToken
	// @Failure 400 {object} server.ErrorResponse
	// @Failure 401 {object} server.ErrorResponse
	// @Failure 500 {object} server.ErrorResponse
	// @Router /auth/refresh-token [post]
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
