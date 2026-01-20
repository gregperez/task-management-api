package domain

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrUserAlreadyExists  = errors.New("el usuario ya existe")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUnauthorized       = errors.New("no autorizado")
	ErrForbidden          = errors.New("acceso prohibido")
	ErrCannotCreateAdmin  = errors.New("no puede crear usuarios administradores")
	ErrMustChangePassword = errors.New("debe cambiar su contraseña temporal")

	// Task errors
	ErrTaskNotFound            = errors.New("tarea no encontrada")
	ErrTaskExpired             = errors.New("la tarea está vencida")
	ErrInvalidStatusTransition = errors.New("transición de estado no válida")
	ErrTaskCannotBeModified    = errors.New("la tarea no puede ser modificada en su estado actual")
	ErrNotAssignedToUser       = errors.New("tarea no asignada a este usuario")

	// Validation errors
	ErrInvalidInput = errors.New("entrada inválida")
	ErrEmptyField   = errors.New("campo requerido vacío")
)
