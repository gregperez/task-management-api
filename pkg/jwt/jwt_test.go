package jwt

import (
	"testing"
	"time"
)

func TestJwt_GenerateTokenAndValidateToken(t *testing.T) {
	secret := "mysecret"
	userID := "12345"
	role := "admin"
	expiry := time.Minute * 15
	
	token, err := GenerateToken(userID, role, secret, expiry)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}