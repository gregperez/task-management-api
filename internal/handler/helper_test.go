package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/middleware"
)

func TestHelper_GetUserID(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantUserID string
		wantErr   error
	}{
		{
			name:      "Valid UserID in context",
			ctx:       context.WithValue(context.Background(), middleware.UserIDKey, "user-123"),
			wantUserID: "user-123",
			wantErr:   nil,
		},
		{
			name:      "Missing UserID in context",
			ctx:       context.Background(),
			wantUserID: "",
			wantErr:   ErrMissingUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := GetUserID(tt.ctx)
			if gotUserID != tt.wantUserID {
				t.Errorf("GetUserID() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if err != tt.wantErr {
				t.Errorf("GetUserID() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelper_GetUserRole(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		wantUserRole domain.UserRole
		wantErr     error
	}{
		{
			name:        "Valid UserRole Admin in context",
			ctx:         context.WithValue(context.Background(), middleware.UserRoleKey, domain.RoleAdmin),
			wantUserRole: domain.RoleAdmin,
			wantErr:     nil,
		},
		{
			name:        "Missing UserRole in context",
			ctx:         context.Background(),
			wantUserRole: "",
			wantErr:     ErrMissingUserRole,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserRole, err := GetUserRole(tt.ctx)
			if gotUserRole != tt.wantUserRole {
				t.Errorf("GetUserRole() gotUserRole = %v, want %v", gotUserRole, tt.wantUserRole)
			}
			if err != tt.wantErr {
				t.Errorf("GetUserRole() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelper_GetUserContext(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		wantUserID    string
		wantUserRole  domain.UserRole
		wantErr      error
	}{
		{
			name:         "Valid UserID and UserRole in context",
			ctx:          context.WithValue(context.WithValue(context.Background(), middleware.UserIDKey, "user-123"), middleware.UserRoleKey, domain.RoleExecutor),
			wantUserID:    "user-123",
			wantUserRole:  domain.RoleExecutor,
			wantErr:      nil,
		},
		{
			name:         "Missing UserID in context",
			ctx:          context.WithValue(context.Background(), middleware.UserRoleKey, domain.RoleExecutor),
			wantUserID:    "",
			wantUserRole:  "",
			wantErr:      ErrMissingUserID,
		},
		{
			name:         "Missing UserRole in context",
			ctx:          context.WithValue(context.Background(), middleware.UserIDKey, "user-123"),
			wantUserID:    "",
			wantUserRole:  "",
			wantErr:      ErrMissingUserRole,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, gotUserRole, err := GetUserContext(tt.ctx)
			if gotUserID != tt.wantUserID {
				t.Errorf("GetUserContext() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if gotUserRole != tt.wantUserRole {
				t.Errorf("GetUserContext() gotUserRole = %v, want %v", gotUserRole, tt.wantUserRole)
			}
			if err != tt.wantErr {
				t.Errorf("GetUserContext() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelper_DecodeJSON(t *testing.T) {
	type SampleRequest struct {
		UserName  string `json:"username"`
		Email string `json:"email"`
	}
	tests := []struct {
		name       string
		jsonInput  string
		wantStruct SampleRequest
		wantErr    error
	}{
		{
			name:       "Valid JSON input",
			jsonInput:  `{"username": "John Doe", "email": "john.doe@example.com"}`,
			wantStruct: SampleRequest{UserName: "John Doe", Email: "john.doe@example.com"},
			wantErr:    nil,
		},
		{
			name:       "Invalid JSON input",
			jsonInput:  `{"username": "John Doe", "email": "`,
			wantStruct: SampleRequest{},
			wantErr:    ErrInvalidJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := http.Request{
				Body: io.NopCloser(strings.NewReader(tt.jsonInput)),
			}
			var gotStruct SampleRequest
			err := DecodeJSON(&r, &gotStruct)
			if gotStruct != tt.wantStruct {
				t.Errorf("DecodeJSON() gotStruct = %v, want %v", gotStruct, tt.wantStruct)
			}
			if err != tt.wantErr {
				t.Errorf("DecodeJSON() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelper_HTTPStatusFromError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantCode int
	}{
		{
			name:    "No error",
			err:     nil,
			wantCode: http.StatusOK,
		},
		{
			name:    "Domain NotFound error",
			err:     domain.ErrUserNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:    "Domain Unauthorized error",
			err:     domain.ErrUnauthorized,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:    "Domain Forbidden error",
			err:     domain.ErrForbidden,
			wantCode: http.StatusForbidden,
		},
		{
			name:    "Domain InvalidInput error",
			err:     domain.ErrInvalidInput,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "Domain UserAlreadyExists error",
			err:     domain.ErrUserAlreadyExists,
			wantCode: http.StatusConflict,
		},
		{
			name:    "Handler BadRequest error",
			err:     ErrInvalidJSON,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "Handler StatusUnauthorized error",
			err:     ErrMissingUserID,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:    "Unknown error",
			err:     errors.New("some unknown error"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode := MapErrorToStatus(tt.err)
			if gotCode != tt.wantCode {
				t.Errorf("MapErrorToStatus() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
