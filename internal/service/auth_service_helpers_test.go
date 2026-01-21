package service

import (
	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/pkg/password"
	"testing"
	"time"
)

func TestGenerateAuthToken(t *testing.T) {
	jwtSecret := "testsecret"
	tokenExpiry := time.Minute * 30
	userID := "user123"
	role := domain.RoleAdmin

	token, err := generateAuthToken(userID, role, jwtSecret, tokenExpiry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("Expected a token, got empty string")
	}
}

func TestBuildLoginResponse(t *testing.T) {
	user := &domain.User{
		ID:                 "user123",
		Role:               domain.RoleAuditor,
		IsTemporaryPassword:  false,
	}
	token := "testtoken"
	response := buildLoginResponse(user, token)
	if response.Token != token {
		t.Errorf("Expected token %s, got %s", token, response.Token)
	}
	if response.User != user {
		t.Errorf("Expected user %v, got %v", user, response.User)
	}
	if response.MustChangePassword != false {
		t.Errorf("Expected MustChangePassword false, got %v", response.MustChangePassword)
	}
}

func TestUpdatePasswordFields(t *testing.T) {
	user := &domain.User{
		Password:            "oldhash",
		IsTemporaryPassword:  true,
		UpdatedAt:           time.Now().Add(-time.Hour),
	}
	newPasswordHash := "newhash"
	updatePasswordFields(user, newPasswordHash)
	if user.Password != newPasswordHash {
		t.Errorf("Expected password %s, got %s", newPasswordHash, user.Password)
	}
	if user.IsTemporaryPassword != false {
		t.Errorf("Expected IsTemporaryPassword false, got %v", user.IsTemporaryPassword)
	}
	if user.UpdatedAt.Before(time.Now().Add(-time.Minute)) {
		t.Errorf("Expected UpdatedAt to be recent, got %v", user.UpdatedAt)
	}
}

func TestValidateCredentials(t *testing.T) {
	hashedPassword, _ := password.HashPassword("correctpassword")
	err := validateCredentials("correctpassword", hashedPassword)
	if err != nil {
		t.Errorf("Expected no error for correct password, got %v", err)
	}
	err = validateCredentials("wrongpassword", hashedPassword)
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials for wrong password, got %v", err)
	}
}