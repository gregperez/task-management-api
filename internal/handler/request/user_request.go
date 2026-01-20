package request

// UserRequests representa las estructuras de solicitud para operaciones de usuarios

// CreateUserRequest representa la solicitud para crear un nuevo usuario
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
}
