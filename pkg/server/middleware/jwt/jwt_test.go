package jwt

import (
	"reflect"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func TestService_MWFunc(t *testing.T) {
	type fields struct {
		Algo            jwt.SigningMethod
		SecretKey       []byte
		AccessDuration  time.Duration
		RefreshDuration time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   echo.MiddlewareFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Service{
				Algo:            tt.fields.Algo,
				SecretKey:       tt.fields.SecretKey,
				AccessDuration:  tt.fields.AccessDuration,
				RefreshDuration: tt.fields.RefreshDuration,
			}
			if got := j.MWFunc(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.MWFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
