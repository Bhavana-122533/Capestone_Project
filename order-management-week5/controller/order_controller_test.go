package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"order-management/authn"
	"order-management/repository"
	"order-management/service"

	"github.com/gin-gonic/gin"
)

func TestParseID(t *testing.T) {
	id, err := parseID("42")
	if err != nil {
		t.Fatalf("parseID() returned error: %v", err)
	}
	if id != 42 {
		t.Fatalf("parseID() = %d, want 42", id)
	}

	if _, err := parseID("abc"); err == nil {
		t.Fatal("parseID() expected error for invalid value")
	}
}

func TestParseIntQuery(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/orders?limit=10&offset=3", nil)

	limit, err := parseIntQuery(c, "limit", 20)
	if err != nil || limit != 10 {
		t.Fatalf("parseIntQuery(limit) = (%d, %v), want (10, nil)", limit, err)
	}

	offset, err := parseIntQuery(c, "offset", 0)
	if err != nil || offset != 3 {
		t.Fatalf("parseIntQuery(offset) = (%d, %v), want (3, nil)", offset, err)
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/orders?offset=1", nil)
	fallback, err := parseIntQuery(c2, "limit", 20)
	if err != nil || fallback != 20 {
		t.Fatalf("parseIntQuery(fallback) = (%d, %v), want (20, nil)", fallback, err)
	}

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest(http.MethodGet, "/orders", nil)
	fallback, err = parseIntQuery(c3, "limit", 20)
	if err != nil || fallback != 20 {
		t.Fatalf("parseIntQuery(fallback) = (%d, %v), want (20, nil)", fallback, err)
	}
}

func TestWriteOrderError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
	}{
		{name: "forbidden", err: authn.ErrForbidden, statusCode: http.StatusForbidden},
		{name: "invalid", err: service.ErrInvalid, statusCode: http.StatusBadRequest},
		{name: "product unavailable", err: service.ErrProductUnavailable, statusCode: http.StatusBadGateway},
		{name: "not found", err: repository.ErrNotFound, statusCode: http.StatusNotFound},
		{name: "default", err: assertError("fallback"), statusCode: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/orders/1", nil)
			writeOrderError(c, tt.err)
			if w.Code != tt.statusCode {
				t.Fatalf("writeOrderError(%s) status = %d, want %d", tt.name, w.Code, tt.statusCode)
			}
		})
	}
}

type assertError string

func (e assertError) Error() string { return string(e) }
