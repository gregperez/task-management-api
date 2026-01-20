package domain

import (
	"time"
)

type UserRole string

const (
	RoleAdmin    UserRole = "Administrador"
	RoleExecutor UserRole = "Ejecutor"
	RoleAuditor  UserRole = "Auditor"
)

type User struct {
	ID                  string    `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Password            string    `json:"-"` // No se serializa en JSON
	Role                UserRole  `json:"role"`
	IsTemporaryPassword bool      `json:"is_temporary_password"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsExecutor() bool {
	return u.Role == RoleExecutor
}

func (u *User) IsAuditor() bool {
	return u.Role == RoleAuditor
}

func (u *User) MustChangePassword() bool {
	return u.IsTemporaryPassword
}
