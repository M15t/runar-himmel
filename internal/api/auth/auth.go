package auth

import (
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/server"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// Login tries to authenticate the user provided by given credentials
func (s *Auth) Login(c echo.Context, data Credentials) (*types.AuthToken, error) {
	existedUser, err := s.repo.User.FindByEmail(c.Request().Context(), data.Email)
	if err != nil || existedUser == nil {
		return nil, ErrInvalidCredentials.SetInternal(err)
	}

	if !s.cr.CompareHashAndPassword(existedUser.Password, data.Password) {
		return nil, ErrInvalidCredentials
	}

	switch data.GrantType {
	case "app":
		if existedUser.Role != rbac.RoleUser {
			return nil, ErrInvalidCredentials
		}
	case "portal":
		if !lo.Contains([]string{rbac.RoleAdmin, rbac.RoleSuperAdmin}, existedUser.Role) {
			return nil, ErrInvalidCredentials
		}
	default:
		return nil, ErrInvalidGrantType
	}

	if existedUser.Status == types.UserStatusBlocked.String() {
		return nil, ErrUserBlocked
	}

	return s.authenticate(c, &AuthenticateInput{
		User:    existedUser,
		IsLogin: true,
	})
}

// RefreshToken refreshes the access token
func (s *Auth) RefreshToken(c echo.Context, data RefreshTokenData) (*types.AuthToken, error) {
	token, err := s.jwt.ParseToken(data.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken.SetInternal(err)
	}

	// claims token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	// get user id and session id from claims
	sessionID, ok := claims["id"].(string)
	userID, ok := claims["uid"].(string)
	existedSession, err := s.repo.Session.FindByID(c.Request().Context(), sessionID, userID)
	if err != nil || existedSession == nil {
		return nil, ErrInvalidRefreshToken.SetInternal(err)
	}

	// check if session is expired
	if time.Now().After(existedSession.ExpiresAt) {
		// update session to blocked
		if err := s.repo.Session.Update(c.Request().Context(), &types.Session{
			IsBlocked: true,
		}, existedSession.ID); err != nil {
			return nil, ErrInvalidRefreshToken.SetInternal(err)
		}

		return nil, ErrTokenExpired
	}

	// update session
	if err := s.repo.Session.Update(c.Request().Context(), &types.Session{
		IPAddress: c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}, sessionID); err != nil {
		return nil, server.NewHTTPInternalError("error updating session").SetInternal(err)
	}

	return s.authenticate(c, &AuthenticateInput{
		User:    existedSession.User,
		IsLogin: false,
	})
}
