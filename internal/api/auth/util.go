package auth

import (
	"context"
	"fmt"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/server/middleware/jwt"
)

func (s *Auth) authenticate(ctx context.Context, u *types.User) (*types.AuthToken, error) {
	accessTokenOutput := jwt.TokenOutput{}
	refreshTokenOutput := jwt.TokenOutput{}
	if err := s.jwt.GenerateToken(&jwt.TokenInput{
		Type: jwt.TypeTokenAccess,
		Claims: map[string]interface{}{
			"id":    u.ID,
			"email": u.Email,
			"name":  fmt.Sprintf("%s %s", u.FirstName, u.LastName),
			"role":  u.Role,
		},
	}, &accessTokenOutput); err != nil {
		return nil, err
	}

	if err := s.jwt.GenerateToken(&jwt.TokenInput{
		Type: jwt.TypeTokenRefresh,
		Claims: map[string]interface{}{
			"id": u.ID,
		},
	}, &refreshTokenOutput); err != nil {
		return nil, err
	}

	// store refresh token in db
	if err := s.repo.User.UpdateRefreshToken(ctx, u.ID, refreshTokenOutput.Token); err != nil {
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
