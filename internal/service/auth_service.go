package service

import (
	"context"
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
	"gregperez/task-management-api/pkg/jwt"
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

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token              string       `json:"token"`
	User               *domain.User `json:"user"`
	MustChangePassword bool         `json:"must_change_password"`
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Verificar contraseña
	if !password.CheckPassword(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generar token JWT
	token, err := jwt.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              token,
		User:               user,
		MustChangePassword: user.MustChangePassword(),
	}, nil
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar contraseña actual
	if !password.CheckPassword(req.CurrentPassword, user.Password) {
		return domain.ErrInvalidCredentials
	}

	// Hash nueva contraseña
	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.IsTemporaryPassword = false
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}
