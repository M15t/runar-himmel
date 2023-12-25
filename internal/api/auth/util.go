package auth

import (
	"fmt"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/server/middleware/jwt"
	"runar-himmel/pkg/util/ulidutil"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Auth) authenticate(c echo.Context, ai *AuthenticateInput) (*types.AuthToken, error) {
	// * noted: this case is support only 1 session per user
	// * multiple sessions per user is not supported yet need to be implemented later
	ctx := c.Request().Context()

	sessionID := ulidutil.NewString()
	accessTokenOutput := jwt.TokenOutput{}
	refreshTokenOutput := jwt.TokenOutput{}
	if err := s.jwt.GenerateToken(&jwt.TokenInput{
		Type: jwt.TypeTokenAccess,
		Claims: map[string]interface{}{
			"id":    ai.User.ID,
			"email": ai.User.Email,
			"name":  fmt.Sprintf("%s %s", ai.User.FirstName, ai.User.LastName),
			"role":  ai.User.Role,
		},
	}, &accessTokenOutput); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"last_login": time.Now(),
	}
	if ai.IsLogin {
		if err := s.jwt.GenerateToken(&jwt.TokenInput{
			Type: jwt.TypeTokenRefresh,
			Claims: map[string]interface{}{
				"id":  sessionID,
				"uid": ai.User.ID,
			},
		}, &refreshTokenOutput); err != nil {
			return nil, err
		}

		// create session
		if err := s.repo.Session.Create(ctx, &types.Session{
			ID:        sessionID,
			UserID:    ai.User.ID,
			IPAddress: c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			ExpiresAt: time.Now().Add(time.Duration(refreshTokenOutput.ExpiresIn) * time.Second),
		}); err != nil {
			return nil, err
		}

		updates["refresh_token"] = refreshTokenOutput.Token
	}

	// update last_login and refresh_token
	if err := s.repo.User.Update(ctx, updates, ai.User.ID); err != nil {
		return nil, err
	}

	// TODO: add more logic if needed

	return &types.AuthToken{
		AccessToken:  accessTokenOutput.Token,
		TokenType:    "bearer",
		ExpiresIn:    accessTokenOutput.ExpiresIn,
		RefreshToken: refreshTokenOutput.Token,
	}, nil
}
