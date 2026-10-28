package main

import (
	"testing"

	"order-management/config"
	"order-management/servicehost"
)

func TestAuthServiceUsesExpectedPortFallback(t *testing.T) {
	_, port, err := servicehost.Port("", "8083")
	if err != nil {
		t.Fatalf("servicehost.Port() returned error: %v", err)
	}
	if port != 8083 {
		t.Fatalf("servicehost.Port() = %d, want 8083", port)
	}
}

func TestAuthServiceConfigLoadsBootstrapCredentials(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_PORT", "8083")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("BOOTSTRAP_ADMIN_USERNAME", "admin")
	t.Setenv("BOOTSTRAP_ADMIN_PASSWORD", "SecureAdminPassword12345")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() returned error: %v", err)
	}
	if cfg.BootstrapAdminUsername != "admin" {
		t.Fatalf("BootstrapAdminUsername = %q, want admin", cfg.BootstrapAdminUsername)
	}
	if cfg.BootstrapAdminPassword != "SecureAdminPassword12345" {
		t.Fatalf("BootstrapAdminPassword not loaded as expected")
	}
}
