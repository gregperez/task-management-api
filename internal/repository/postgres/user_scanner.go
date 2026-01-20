package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"gregperez/task-management-api/internal/domain"
)

// scanUser escanea una fila de la base de datos en una estructura User
// Centraliza la lógica de scan para evitar duplicación
func scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsTemporaryPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, mapScanError(err)
	}

	return &user, nil
}

// scanUsers escanea múltiples filas en un slice de usuarios
func scanUsers(rows pgx.Rows) ([]*domain.User, error) {
	var users []*domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsTemporaryPassword,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// mapScanError mapea errores de scan a errores del dominio
func mapScanError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	return err
}

// mapDBError mapea errores de base de datos a errores del dominio
func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	// Verificar si es un error de PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			// Usuario ya existe (puede ser username o email)
			return domain.ErrUserAlreadyExists
		case "23503": // foreign_key_violation
			return domain.ErrInvalidInput
		case "23502": // not_null_violation
			return domain.ErrInvalidInput
		}
	}

	return err
}
