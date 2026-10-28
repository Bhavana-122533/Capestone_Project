package servicehost

import (
	"context"
	"testing"

	"order-management/config"
)

func TestPortUsesFallbackWhenConfiguredEmpty(t *testing.T) {
	port, portNo, err := Port("", "8081")
	if err != nil {
		t.Fatalf("Port() returned error: %v", err)
	}
	if port != "8081" {
		t.Fatalf("Port() = %q, want 8081", port)
	}
	if portNo != 8081 {
		t.Fatalf("portNo = %d, want 8081", portNo)
	}
}

func TestPortRejectsInvalidValue(t *testing.T) {
	_, _, err := Port("abc", "8081")
	if err == nil {
		t.Fatal("Port() expected error for invalid value")
	}
}

func TestRegisterReturnsNilForNilClient(t *testing.T) {
	id, err := Register(context.Background(), nil, config.Config{}, "demo-service", 8080, "/health")
	if err != nil {
		t.Fatalf("Register() returned unexpected error: %v", err)
	}
	if id != "" {
		t.Fatalf("Register() = %q, want empty for nil client", id)
	}
}
