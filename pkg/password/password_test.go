package password

import "testing"

func TestPassword_CompareHash(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		expectedResult bool
	}{
		{
			name:           "Correct password",
			password:       "SecurePass123!",
			expectedResult: true,
		},
		{
			name:           "Incorrect password",
			password:       "WrongPass456!",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword("SecurePass123!")
			if err != nil {
				t.Fatalf("Error hashing password: %v", err)
			}
			result := CheckPassword(tt.password, hashedPassword)
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestGenerateTemporaryPassword(t *testing.T) {
	password := GenerateTemporaryPassword()
	if len(password) != 12 {
		t.Errorf("Expected password length of 12, got %d", len(password))
	}
}