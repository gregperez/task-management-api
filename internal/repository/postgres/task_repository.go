package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gregperez/task-management-api/internal/domain"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, query,
		task.ID, task.Title, task.Description, task.Status,
		task.DueDate, task.AssignedTo, task.CreatedBy, task.CreatedAt, task.UpdatedAt,
	)

	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task domain.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status,
		&task.DueDate, &task.AssignedTo, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	// Cargar comentarios
	comments, err := r.getCommentsByTaskID(ctx, id)
	if err != nil {
		return nil, err
	}
	task.Comments = comments

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query,
		task.Title, task.Description, task.Status, task.DueDate, task.UpdatedAt, task.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) ListByAssignee(ctx context.Context, userID string) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at
		FROM tasks
		WHERE assigned_to = $1
		ORDER BY due_date ASC
	`

	return r.queryTasks(ctx, query, userID)
}

func (r *TaskRepository) ListAll(ctx context.Context) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	return r.queryTasks(ctx, query)
}

func (r *TaskRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, task_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		comment.ID, comment.TaskID, comment.UserID, comment.Content, comment.CreatedAt,
	)

	return err
}

func (r *TaskRepository) queryTasks(ctx context.Context, query string, args ...interface{}) ([]*domain.Task, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.DueDate, &task.AssignedTo, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Cargar comentarios para cada tarea
		comments, err := r.getCommentsByTaskID(ctx, task.ID)
		if err != nil {
			return nil, err
		}
		task.Comments = comments

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *TaskRepository) getCommentsByTaskID(ctx context.Context, taskID string) ([]domain.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
