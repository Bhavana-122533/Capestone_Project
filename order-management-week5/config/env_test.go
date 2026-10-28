package config

import (
	"testing"
	"time"
)

func TestLoadReadsEnvironmentValues(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "orders")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("CONSUL_HTTP_ADDR", "http://localhost:8500")
	t.Setenv("SERVICE_NAME", "order-service")
	t.Setenv("SERVICE_ID", "order-1")
	t.Setenv("SERVICE_ADDRESS", "127.0.0.1")
	t.Setenv("PRODUCT_SERVICE_URL", "http://localhost:4042")
	t.Setenv("PRODUCT_TIMEOUT_MS", "1500")
	t.Setenv("PRODUCT_RETRIES", "3")
	t.Setenv("PRODUCT_BOOTSTRAP_ID", "501")
	t.Setenv("PRODUCT_BOOTSTRAP_NAME", "Sample Product")
	t.Setenv("PRODUCT_BOOTSTRAP_PRICE", "12.50")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.AppPort != "9090" {
		t.Fatalf("AppPort = %q, want 9090", cfg.AppPort)
	}
	if cfg.JWTSecret != "12345678901234567890123456789012" {
		t.Fatalf("JWTSecret not loaded correctly")
	}
	if cfg.ProductTimeout != 1500*time.Millisecond {
		t.Fatalf("ProductTimeout = %v, want 1500ms", cfg.ProductTimeout)
	}
	if cfg.ProductRetries != 3 {
		t.Fatalf("ProductRetries = %d, want 3", cfg.ProductRetries)
	}
	if cfg.ProductBootstrapPrice != 12.50 {
		t.Fatalf("ProductBootstrapPrice = %v, want 12.5", cfg.ProductBootstrapPrice)
	}
}

func TestLoadRequiresJWTSecretInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "short")
	t.Setenv("CONSUL_HTTP_ADDR", "http://localhost:8500")
	t.Setenv("SERVICE_ADDRESS", "127.0.0.1")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error in production when JWT_SECRET is too short")
	}
}
