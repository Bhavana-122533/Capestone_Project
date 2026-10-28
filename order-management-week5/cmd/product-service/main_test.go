package main

import (
	"testing"

	"order-management/config"
	"order-management/servicehost"
)

func TestProductServiceUsesExpectedPortFallback(t *testing.T) {
	_, port, err := servicehost.Port("", "8082")
	if err != nil {
		t.Fatalf("servicehost.Port() returned error: %v", err)
	}
	if port != 8082 {
		t.Fatalf("servicehost.Port() = %d, want 8082", port)
	}
}

func TestProductServiceConfigLoadsBootstrapProduct(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_PORT", "4042")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("PRODUCT_BOOTSTRAP_ID", "501")
	t.Setenv("PRODUCT_BOOTSTRAP_NAME", "Sample Product")
	t.Setenv("PRODUCT_BOOTSTRAP_PRICE", "12.50")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() returned error: %v", err)
	}
	if cfg.ProductBootstrapID != 501 {
		t.Fatalf("ProductBootstrapID = %d, want 501", cfg.ProductBootstrapID)
	}
	if cfg.ProductBootstrapName != "Sample Product" {
		t.Fatalf("ProductBootstrapName = %q, want Sample Product", cfg.ProductBootstrapName)
	}
	if cfg.ProductBootstrapPrice != 12.5 {
		t.Fatalf("ProductBootstrapPrice = %v, want 12.5", cfg.ProductBootstrapPrice)
	}
}
