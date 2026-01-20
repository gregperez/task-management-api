package middleware

import (
	"errors"
	"net/http"
	"testing"

	"gregperez/task-management-api/internal/domain"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    error
	}{
		{
			name:       "Valid token",
			authHeader: "Bearer validtoken123",
			wantToken:  "validtoken123",
			wantErr:    nil,
		},
		{
			name:       "Missing token",
			authHeader: "",
			wantToken:  "",
			wantErr:    ErrMissingToken,
		},
		{
			name:       "Invalid format - no Bearer",
			authHeader: "InvalidTokenFormat",
			wantToken:  "",
			wantErr:    ErrInvalidTokenFormat,
		},
		{
			name:       "Invalid format - empty token",
			authHeader: "Bearer ",
			wantToken:  "",
			wantErr:    ErrInvalidTokenFormat,
		},
		{
			name:       "Invalid format - wrong scheme",
			authHeader: "THEBearer ",
			wantToken:  "",
			wantErr:    ErrInvalidTokenFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := ExtractBearerToken(tt.authHeader)
			if gotToken != tt.wantToken {
				t.Errorf("ExtractBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if err != tt.wantErr {
				t.Errorf("ExtractBearerToken() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHasRequiredRole(t *testing.T) {
	tests := []struct {
		name         string
		userRole     domain.UserRole
		allowedRoles []domain.UserRole
		want         bool
	}{
		{
			name:         "Role is allowed",
			userRole:     domain.RoleAdmin,
			allowedRoles: []domain.UserRole{domain.RoleAdmin, domain.UserRole("user")},
			want:         true,
		},
		{
			name:         "Role is not allowed",
			userRole:     domain.UserRole("guest"),
			allowedRoles: []domain.UserRole{domain.RoleAdmin, domain.UserRole("user")},
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasRequiredRole(tt.userRole, tt.allowedRoles); got != tt.want {
				t.Errorf("HasRequiredRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapJWTError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantMsg  string
		wantCode string
	}{
		{
			name:     "Expired token error",
			err:      errors.New("expired"),
			wantMsg:  MsgExpiredToken,
			wantCode: ErrCodeExpiredToken,
		},
		{
			name:     "Malformed token error",
			err:      errors.New("malformed"),
			wantMsg:  MsgInvalidTokenFormat,
			wantCode: ErrCodeInvalidTokenFormat,
		},
		{
			name:     "Unknown token error",
			err:      errors.New("invalid"),
			wantMsg:  MsgInvalidToken,
			wantCode: ErrCodeInvalidToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotCode := MapJWTError(tt.err)
			if gotMsg != tt.wantMsg {
				t.Errorf("MapJWTError() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotCode != tt.wantCode {
				t.Errorf("MapJWTError() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

type mockResponseWriter struct {
	header     http.Header
	statusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func TestRespondAuthError(t *testing.T) {
	rr := &mockResponseWriter{header: make(http.Header)}
	RespondAuthError(rr, 401, "Test message", "TEST_CODE")
	if rr.statusCode != 401 {
		t.Errorf("RespondAuthError() statusCode = %v, want %v", rr.statusCode, 401)
	}
}