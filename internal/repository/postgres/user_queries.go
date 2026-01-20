package postgres

// Constantes de columnas para evitar duplicación
const (
	userColumns = "id, username, email, password, role, is_temporary_password, created_at, updated_at"
)

// Queries SQL para operaciones de usuarios
const (
	queryInsertUser = `
		INSERT INTO users (id, username, email, password, role, is_temporary_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	querySelectUserByID = `
		SELECT ` + userColumns + `
		FROM users
		WHERE id = $1
	`

	querySelectUserByUsername = `
		SELECT ` + userColumns + `
		FROM users
		WHERE username = $1
	`

	querySelectUserByEmail = `
		SELECT ` + userColumns + `
		FROM users
		WHERE email = $1
	`

	queryUpdateUser = `
		UPDATE users
		SET email = $1, password = $2, role = $3, is_temporary_password = $4, updated_at = $5
		WHERE id = $6
	`

	queryDeleteUser = `
		DELETE FROM users WHERE id = $1
	`

	queryListUsers = `
		SELECT ` + userColumns + `
		FROM users
		ORDER BY created_at DESC
	`
)
