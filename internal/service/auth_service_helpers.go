package service

import (
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/service/dto"
	"gregperez/task-management-api/pkg/jwt"
	"gregperez/task-management-api/pkg/password"
)

// generateAuthToken genera un token JWT para el usuario.
func generateAuthToken(userID string, role domain.UserRole, jwtSecret string, tokenExpiry time.Duration) (string, error) {
	return jwt.GenerateToken(userID, string(role), jwtSecret, tokenExpiry)
}

// buildLoginResponse construye la respuesta de login con el token y datos del usuario.
func buildLoginResponse(user *domain.User, token string) *dto.LoginResponse {
	return &dto.LoginResponse{
		Token:              token,
		User:               user,
		MustChangePassword: user.MustChangePassword(),
	}
}

// updatePasswordFields actualiza la contraseña del usuario y marca que no es temporal.
func updatePasswordFields(user *domain.User, newPasswordHash string) {
	user.Password = newPasswordHash
	user.IsTemporaryPassword = false
	user.UpdatedAt = now()
}

// validateCredentials verifica que las credenciales sean correctas.
// Retorna ErrInvalidCredentials si la contraseña no coincide.
func validateCredentials(plainPassword, hashedPassword string) error {
	if !password.CheckPassword(plainPassword, hashedPassword) {
		return domain.ErrInvalidCredentials
	}
	return nil
}
