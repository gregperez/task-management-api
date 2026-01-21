package request

// AuthRequests representa las estructuras de solicitud para operaciones de autenticación

// LoginRequest representa la solicitud de inicio de sesión
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword 	string `json:"new_password"`
}
