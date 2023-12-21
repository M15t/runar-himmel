package server_test

import (
	"testing"

	"runar-himmel/pkg/server"
)

// Improve tests
func TestNew(t *testing.T) {
	cfg := &server.Config{Port: 8080}

	e := server.New(cfg)
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
