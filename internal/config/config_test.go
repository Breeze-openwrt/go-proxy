package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("../../config.jsonc")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ListenAddr != ":8080" {
		t.Errorf("Expected addr :8080, got %s", cfg.ListenAddr)
	}

	route, ok := cfg.Routes["example.com"]
	if !ok {
		t.Fatal("Route example.com not found")
	}

	if route.JumpStart != 4 {
		t.Errorf("Expected jump_start 4, got %d", route.JumpStart)
	}
}
