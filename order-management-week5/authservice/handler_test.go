package authservice

import (
	"strings"
	"testing"
)

func TestBootstrapAdminRejectsInvalidUsernameOrPassword(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{name: "short username", username: "ab", password: "StrongPassword123!"},
		{name: "short password", username: "admin", password: "short"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BootstrapAdmin(nil, tt.username, tt.password)
			if err == nil {
				t.Fatalf("BootstrapAdmin() expected error for %q", tt.name)
			}
			if !strings.Contains(err.Error(), "BOOTSTRAP_ADMIN") {
				t.Fatalf("BootstrapAdmin() error = %q, want bootstrap validation message", err.Error())
			}
		})
	}
}
