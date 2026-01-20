package service

import (
	"context"
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
	"gregperez/task-management-api/pkg/password"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type CreateUserRequest struct {
	Username string          `json:"username" validate:"required,min=3"`
	Email    string          `json:"email" validate:"required,email"`
	Role     domain.UserRole `json:"role" validate:"required,oneof=Ejecutor Auditor"`
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest, creatorRole domain.UserRole) (*domain.User, string, error) {
	// Solo administradores pueden crear usuarios
	if creatorRole != domain.RoleAdmin {
		return nil, "", domain.ErrForbidden
	}

	// No se pueden crear administradores
	if req.Role == domain.RoleAdmin {
		return nil, "", domain.ErrCannotCreateAdmin
	}

	// Verificar que el usuario no exista
	existing, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existing != nil {
		return nil, "", domain.ErrUserAlreadyExists
	}

	// Generar contraseña temporal
	tempPassword := password.GenerateTemporaryPassword()
	hashedPassword, err := password.HashPassword(tempPassword)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		ID:                  uuid.New().String(),
		Username:            req.Username,
		Email:               req.Email,
		Password:            hashedPassword,
		Role:                req.Role,
		IsTemporaryPassword: true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	return user, tempPassword, nil
}

type UpdateUserRequest struct {
	Email string          `json:"email" validate:"omitempty,email"`
	Role  domain.UserRole `json:"role" validate:"omitempty,oneof=Ejecutor Auditor"`
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest, updaterRole domain.UserRole) error {
	if updaterRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Role != "" {
		if req.Role == domain.RoleAdmin {
			return domain.ErrCannotCreateAdmin
		}
		user.Role = req.Role
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, userID string, deleterRole domain.UserRole) error {
	if deleterRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	return s.userRepo.Delete(ctx, userID)
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) ListUsers(ctx context.Context, requesterRole domain.UserRole) ([]*domain.User, error) {
	if requesterRole != domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}

	return s.userRepo.List(ctx)
}
