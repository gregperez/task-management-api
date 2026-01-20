package service

import (
	"context"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
	"gregperez/task-management-api/internal/service/dto"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser crea un nuevo usuario con contraseña temporal.
// Solo los administradores pueden crear usuarios.
// No se permite crear usuarios con rol Admin.
func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest, creatorRole domain.UserRole) (*domain.User, string, error) {
	if err := canCreateUser(creatorRole); err != nil {
		return nil, "", err
	}

	if err := validateNotAdminRole(req.Role); err != nil {
		return nil, "", err
	}

	if err := validateUserNotExists(ctx, s.userRepo, req.Username); err != nil {
		return nil, "", err
	}

	user, tempPassword, err := buildNewUser(req)
	if err != nil {
		return nil, "", err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	return user, tempPassword, nil
}

// UpdateUser actualiza los campos de un usuario existente.
// Solo los administradores pueden actualizar usuarios.
// No se permite cambiar el rol a Admin.
func (s *UserService) UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest, updaterRole domain.UserRole) error {
	if err := canUpdateUser(updaterRole); err != nil {
		return err
	}

	if err := validateRoleChange(req.Role); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	updateUserFields(user, req)
	markUserAsUpdated(user)

	return s.userRepo.Update(ctx, user)
}

// DeleteUser elimina un usuario del sistema.
// Solo los administradores pueden eliminar usuarios.
func (s *UserService) DeleteUser(ctx context.Context, userID string, deleterRole domain.UserRole) error {
	if err := canDeleteUser(deleterRole); err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, userID)
}

// GetUser obtiene un usuario por su ID.
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// ListUsers retorna todos los usuarios del sistema.
// Solo los administradores pueden listar usuarios.
func (s *UserService) ListUsers(ctx context.Context, requesterRole domain.UserRole) ([]*domain.User, error) {
	if err := canListUsers(requesterRole); err != nil {
		return nil, err
	}

	return s.userRepo.List(ctx)
}
