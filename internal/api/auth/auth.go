package auth

import (
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/types"

	"github.com/labstack/echo/v4"
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
		if existedUser.Role != rbac.RoleCustomer {
			return nil, ErrInvalidCredentials
		}
	case "portal":
		if existedUser.Role != rbac.RoleAdmin {
			return nil, ErrInvalidCredentials
		}
	default:
		return nil, ErrInvalidGrantType
	}

	if existedUser.Status == types.UserStatusBlocked.String() {
		return nil, ErrUserBlocked
	}

	return s.authenticate(c.Request().Context(), existedUser)
}

// RefreshToken refreshes the access token
func (s *Auth) RefreshToken(echo.Context, RefreshTokenData) (*types.AuthToken, error) {
	return nil, nil
}
