package secure

import (
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestBodyDump(t *testing.T) {
	tests := []struct {
		name string
		want echo.MiddlewareFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BodyDump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BodyDump() = %v, want %v", got, tt.want)
			}
		})
	}
}
