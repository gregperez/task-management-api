package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/pkg/jwt"
)

// Helper para crear un token JWT válido para tests
func createTestToken(userID, role, secret string, expiry time.Duration) string {
	token, _ := jwt.GenerateToken(userID, role, secret, expiry)
	return token
}

// Helper para verificar respuesta de error JSON
func verifyErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedCode string) {
	t.Helper()

	if rr.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response AuthErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response.Code != expectedCode {
		t.Errorf("Expected error code %s, got %s", expectedCode, response.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	secret := "testsecret123"
	token := createTestToken("user123", "admin", secret, 1 * time.Hour)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create mock next handler that checks context
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		// Verify user ID in context
		userID, ok := r.Context().Value(UserIDKey).(string)
		if !ok || userID != "user123" {
			t.Error("UserID not found in context or incorrect")
		}

		// Verify role in context
		userRole, ok := r.Context().Value(UserRoleKey).(domain.UserRole)
		if !ok || userRole != "admin" {
			t.Error("UserRole not found in context or incorrect")
		}

		w.WriteHeader(http.StatusOK)
	})

	// Execute middleware
	middleware := Auth(secret)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	// Verify next handler was called
	if !nextCalled {
		t.Error("Next handler was not called")
	}

	// Verify response status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAuth_MissingToken(t *testing.T) {
	secret := "testsecret123"

	// Create request without Authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Create next handler (should not be called)
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	// Execute middleware
	middleware := Auth(secret)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	// Verify next handler was NOT called
	if nextCalled {
		t.Error("Next handler should not be called with missing token")
	}

	// Verify error response
	verifyErrorResponse(t, rr, http.StatusUnauthorized, ErrCodeMissingToken)
}

func TestAuth_InvalidTokenFormat(t *testing.T) {
	secret := "testsecret123"

	tests := []struct {
		name         string
		authHeader   string
		expectedCode string
	}{
		{
			name:         "No Bearer prefix",
			authHeader:   "InvalidToken123",
			expectedCode: ErrCodeInvalidTokenFormat,
		},
		{
			name:         "Empty token",
			authHeader:   "Bearer ",
			expectedCode: ErrCodeInvalidTokenFormat,
		},
		{
			name:         "Wrong scheme",
			authHeader:   "Basic token123",
			expectedCode: ErrCodeInvalidTokenFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rr := httptest.NewRecorder()

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
			})

			middleware := Auth(secret)
			handler := middleware(next)
			handler.ServeHTTP(rr, req)

			if nextCalled {
				t.Error("Next handler should not be called with invalid token format")
			}

			verifyErrorResponse(t, rr, http.StatusUnauthorized, tt.expectedCode)
		})
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	secret := "testsecret123"

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken123")
	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware := Auth(secret)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Error("Next handler should not be called with invalid token")
	}

	// Should get an error (could be invalid token or expired)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	secret := "testsecret123"
	// Create expired token (negative duration)
	token := createTestToken("user123", "admin", secret, -1*time.Hour)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware := Auth(secret)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Error("Next handler should not be called with expired token")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestRequireRole_ValidRole(t *testing.T) {
	// Create request with role in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserRoleKey, domain.RoleAdmin)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Require admin role
	middleware := RequireRole(domain.RoleAdmin)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Error("Next handler should be called when user has required role")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	// Create request with executor role
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserRoleKey, domain.RoleExecutor)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Allow both admin and executor
	middleware := RequireRole(domain.RoleAdmin, domain.RoleExecutor)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Error("Next handler should be called when user has one of the required roles")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestRequireRole_InsufficientPermissions(t *testing.T) {
	// Create request with auditor role
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserRoleKey, domain.RoleAuditor)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	// Require admin role only
	middleware := RequireRole(domain.RoleAdmin)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Error("Next handler should not be called when user lacks required role")
	}

	verifyErrorResponse(t, rr, http.StatusForbidden, ErrCodeInsufficientPermissions)
}

func TestRequireRole_MissingRole(t *testing.T) {
	// Create request without role in context
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware := RequireRole(domain.RoleAdmin)
	handler := middleware(next)
	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Error("Next handler should not be called when role is missing from context")
	}

	verifyErrorResponse(t, rr, http.StatusUnauthorized, ErrCodeMissingRole)
}
