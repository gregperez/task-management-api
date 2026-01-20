package password

import "testing"

func TestPassword_CompareHash(t *testing.T) {
	password := "Admin123!"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	match := CheckPassword(password, hashedPassword)
	if !match {
		t.Fatalf("Password does not match hash")
	}
	wrongPassword := "WrongPassword!"
	match = CheckPassword(wrongPassword, hashedPassword)
	if match {
		t.Fatalf("Wrong password should not match hash")
	}
}
