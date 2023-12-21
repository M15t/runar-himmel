package jwt

import (
	"fmt"
	"net/http"
	"runar-himmel/pkg/server"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// New generates new JWT service necessery for auth middleware
func New(algo, secretKey string, durations ...int) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}

	var accessDuration, refreshDuration int
	switch len(durations) {
	case 0:
		accessDuration = durations[0]
	case 1:
		accessDuration = durations[0]
		refreshDuration = durations[1]
	default:
		// default values
		accessDuration = 1 * 60 * 60   // 1 hour (in seconds)
		refreshDuration = 24 * 60 * 60 // 24 hours (in seconds)
	}

	return &Service{
		Algo:            signingMethod,
		SecretKey:       []byte(secretKey),
		AccessDuration:  time.Duration(accessDuration) * time.Second,
		RefreshDuration: time.Duration(refreshDuration) * time.Second,
	}
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Algo signing algorithm used for signing.
	Algo jwt.SigningMethod
	// SecretKey used for signing.
	SecretKey []byte
	// AccessKeyDuration duration (in seconds) for which the jwt access token is valid.
	AccessDuration time.Duration
	// RefreshDuration duration (in seconds) for which the jwt refresh token is valid.
	RefreshDuration time.Duration
}

// MWFunc makes JWT implement the Middleware interface.
func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseTokenFromHeader(c)
			if err != nil || !token.Valid {
				if err != nil {
					c.Logger().Errorf("error parsing token: %+v", err)
				}
				return server.NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", "Your session is unauthorized or has expired.").SetInternal(err)
			}

			claims := token.Claims.(jwt.MapClaims)
			for key, val := range claims {
				c.Set(key, val)
			}

			return next(c)
		}
	}
}

// ParseTokenFromHeader parses token from Authorization header
func (j *Service) ParseTokenFromHeader(c echo.Context) (*jwt.Token, error) {
	token := c.Request().Header.Get("Authorization")
	if len(strings.TrimSpace(token)) == 0 {
		return nil, fmt.Errorf("token not found")
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
		return nil, fmt.Errorf("token invalid")
	}

	return j.ParseToken(parts[1])
}

// ParseToken parses token from string
func (j *Service) ParseToken(input string) (*jwt.Token, error) {
	return jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
		if j.Algo != token.Method {
			return nil, fmt.Errorf("token method mismatched")
		}
		return j.SecretKey, nil
	})
}

// GenerateToken generates new Service token and populates it with user data
func (j *Service) GenerateToken(input *TokenInput, output *TokenOutput) error {
	if input == nil || output == nil {
		return fmt.Errorf("input and output cannot be nil")
	}

	// Set token expiration based on token type
	var expire time.Time
	switch input.Type {
	case "access_token":
		expire = time.Now().Add(j.AccessDuration)
	case "refresh_token":
		expire = time.Now().Add(j.RefreshDuration)
	default:
		return fmt.Errorf("invalid token type")
	}

	// Set expiration claim
	input.Claims["exp"] = expire.Unix()

	// Create JWT token
	token := jwt.NewWithClaims(j.Algo, jwt.MapClaims(input.Claims))

	// Sign the token
	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %w", err)
	}

	// Populate output struct
	output.Token = tokenString
	output.ExpiresIn = int(expire.Sub(time.Now()).Seconds())

	return nil
}
