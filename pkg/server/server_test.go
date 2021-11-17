package server_test

import (
	"testing"

	"ghoul/pkg/server"
)

// Improve tests
func TestNew(t *testing.T) {
	cfg := &server.Config{Stage: "development", Port: 8080}

	e := server.New(cfg)
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
