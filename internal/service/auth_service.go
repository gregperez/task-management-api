package service

import (
	"context"
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
	"gregperez/task-management-api/internal/service/dto"
	"gregperez/task-management-api/pkg/password"
)

type AuthService struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// Login autentica un usuario y retorna un token JWT.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := validateCredentials(req.Password, user.Password); err != nil {
		return nil, err
	}

	token, err := generateAuthToken(user.ID, user.Role, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return nil, err
	}

	return buildLoginResponse(user, token), nil
}

// ChangePassword cambia la contraseña de un usuario.
// Valida la contraseña actual antes de permitir el cambio.
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := validateCredentials(req.CurrentPassword, user.Password); err != nil {
		return err
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	updatePasswordFields(user, hashedPassword)

	return s.userRepo.Update(ctx, user)
}
