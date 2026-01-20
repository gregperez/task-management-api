package middleware

// Constantes para autenticación
const (
	// BearerScheme es el esquema de autenticación Bearer
	BearerScheme = "Bearer"
	
	// AuthHeaderKey es la clave del header de autorización
	AuthHeaderKey = "Authorization"
)

// Códigos de error para respuestas de autenticación
const (
	ErrCodeMissingToken            = "MISSING_TOKEN"
	ErrCodeInvalidTokenFormat      = "INVALID_TOKEN_FORMAT"
	ErrCodeInvalidToken            = "INVALID_TOKEN"
	ErrCodeExpiredToken            = "EXPIRED_TOKEN"
	ErrCodeMissingRole             = "MISSING_ROLE"
	ErrCodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
)

// Mensajes de error en español
const (
	MsgMissingToken            = "Token de autorización no proporcionado"
	MsgInvalidTokenFormat      = "Formato de token inválido. Se espera: Bearer <token>"
	MsgInvalidToken            = "Token inválido o firma no verificable"
	MsgExpiredToken            = "Token expirado. Por favor, inicie sesión nuevamente"
	MsgMissingRole             = "Rol de usuario no encontrado en el contexto"
	MsgInsufficientPermissions = "Acceso denegado. No tiene permisos suficientes"
)
