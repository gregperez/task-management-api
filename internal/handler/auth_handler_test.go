package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/handler/request"
	"gregperez/task-management-api/internal/middleware"
	"gregperez/task-management-api/internal/service/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService implementa service.AuthServiceInterface para tests
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.LoginResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   request.LoginRequest
		setupMock     func(m *MockAuthService)
		expectedError error
	}{
		{
			name: "Successful login",
			requestBody: request.LoginRequest{
				Username: "testuser",
				Password: "testpass",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, dto.LoginRequest{
					Username: "testuser",
					Password: "testpass",
				}).Return(&dto.LoginResponse{
					Token: "valid-jwt-token",
					User: &domain.User{
						ID:       "user-123",
						Username: "testuser",
						Role:     domain.RoleAdmin,
					},
					MustChangePassword: false,
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid credentials",
			requestBody: request.LoginRequest{
				Username: "wronguser",
				Password: "wrongpass",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, dto.LoginRequest{
					Username: "wronguser",
					Password: "wrongpass",
				}).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedError: domain.ErrInvalidCredentials,
		},
		{
			name: "Bad request - missing fields",
			requestBody: request.LoginRequest{
				Username: "testuser",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, dto.LoginRequest{
					Username: "testuser",
				}).Return(nil, domain.ErrEmptyField)
			},
			expectedError: domain.ErrEmptyField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)
			authHandler := NewAuthHandler(mockService)

			reqBody := `{"username":"` + tt.requestBody.Username + `","password":"` + tt.requestBody.Password + `"}`
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.Login(rr, req)

			switch tt.expectedError {
			case nil:
				assert.Equal(t, http.StatusOK, rr.Code)
			case domain.ErrEmptyField:
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), tt.expectedError.Error())
			case domain.ErrInvalidCredentials:
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	tests := []struct {
		name           		string
		userID         		string
		requestBody    		request.ChangePasswordRequest
		setupMock      		func(m *MockAuthService)
		expectedStatus 		int
		contextWithValueOK 	bool
	}{
		{
			name:   "Successful password change",
			userID: "user-123",
			requestBody: request.ChangePasswordRequest{
				CurrentPassword: "oldpass",
				NewPassword:     "newpass",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ChangePassword", mock.Anything, "user-123", dto.ChangePasswordRequest{
					CurrentPassword: "oldpass",
					NewPassword:     "newpass",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			contextWithValueOK: true,
		},
		{
			name:   "Invalid current password",
			userID: "user-123",
			requestBody: request.ChangePasswordRequest{
				CurrentPassword: "wrongoldpass",
				NewPassword:     "newpass",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ChangePassword", mock.Anything, "user-123", dto.ChangePasswordRequest{
					CurrentPassword: "wrongoldpass",
					NewPassword:     "newpass",
				}).Return(domain.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			contextWithValueOK: true,
		},
		{
			name:   "Error changing password getting user",
			userID: "user-123",
			requestBody: request.ChangePasswordRequest{},
			setupMock: func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			contextWithValueOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)
			authHandler := NewAuthHandler(mockService)

			reqBody := `{"current_password":"` + tt.requestBody.CurrentPassword + `","new_password":"` + tt.requestBody.NewPassword + `"}`
			req := httptest.NewRequest(http.MethodPost, "/auth/change-password", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid-jwt-token")

			// Usar la constante middleware.UserIDKey en lugar de string literal
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
			if !tt.contextWithValueOK {
				ctx = req.Context() // Contexto sin UserID
			}
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			authHandler.ChangePassword(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	// Dado que el logout en JWT stateless se maneja en el cliente,
	// este handler puede simplemente devolver un 200 OK.
	authHandler := NewAuthHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rr := httptest.NewRecorder()
	authHandler.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}