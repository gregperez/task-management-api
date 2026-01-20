package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gregperez/task-management-api/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create crea un nuevo usuario en la base de datos
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(ctx, queryInsertUser,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsTemporaryPassword,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return mapDBError(err)
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, querySelectUserByID, id)
	return scanUser(row)
}

// GetByUsername obtiene un usuario por su nombre de usuario
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, querySelectUserByUsername, username)
	return scanUser(row)
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, querySelectUserByEmail, email)
	return scanUser(row)
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	result, err := r.db.Exec(ctx, queryUpdateUser,
		user.Email,
		user.Password,
		user.Role,
		user.IsTemporaryPassword,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return mapDBError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, queryDeleteUser, id)
	if err != nil {
		return mapDBError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List obtiene todos los usuarios ordenados por fecha de creación
func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.Query(ctx, queryListUsers)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	return scanUsers(rows)
}
