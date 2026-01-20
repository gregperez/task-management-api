package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gregperez/task-management-api/internal/domain"
)

// scanTask escanea una fila de la base de datos en una estructura Task (sin comentarios)
func scanTask(row pgx.Row) (*domain.Task, error) {
	var task domain.Task
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, mapTaskScanError(err)
	}

	return &task, nil
}

// scanTasks escanea múltiples filas en un slice de tareas (sin comentarios)
func scanTasks(rows pgx.Rows) ([]*domain.Task, error) {
	var tasks []*domain.Task

	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// scanComments escanea múltiples filas de comentarios
func scanComments(rows pgx.Rows) ([]domain.Comment, error) {
	var comments []domain.Comment

	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// loadTaskComments carga los comentarios de una tarea
func loadTaskComments(ctx context.Context, db *pgxpool.Pool, taskID string) ([]domain.Comment, error) {
	rows, err := db.Query(ctx, querySelectCommentsByTaskID, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanComments(rows)
}

// loadTasksComments carga los comentarios para múltiples tareas
func loadTasksComments(ctx context.Context, db *pgxpool.Pool, tasks []*domain.Task) error {
	for _, task := range tasks {
		comments, err := loadTaskComments(ctx, db, task.ID)
		if err != nil {
			return err
		}
		task.Comments = comments
	}
	return nil
}

// mapTaskScanError mapea errores de scan a errores del dominio
func mapTaskScanError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrTaskNotFound
	}
	return err
}

// mapTaskDBError mapea errores de base de datos a errores del dominio
func mapTaskDBError(err error) error {
	if err == nil {
		return nil
	}

	// Verificar si es un error de PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return domain.ErrInvalidInput
		case "23503": // foreign_key_violation
			return domain.ErrInvalidInput
		case "23502": // not_null_violation
			return domain.ErrInvalidInput
		}
	}

	return err
}
