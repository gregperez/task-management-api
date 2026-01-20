package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gregperez/task-management-api/internal/domain"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create crea una nueva tarea en la base de datos
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := r.db.Exec(ctx, queryInsertTask,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.AssignedTo,
		task.CreatedBy,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return mapTaskDBError(err)
}

// GetByID obtiene una tarea por su ID (incluye comentarios)
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	row := r.db.QueryRow(ctx, querySelectTaskByID, id)
	task, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	// Cargar comentarios
	comments, err := loadTaskComments(ctx, r.db, task.ID)
	if err != nil {
		return nil, err
	}
	task.Comments = comments

	return task, nil
}

// Update actualiza una tarea existente
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	result, err := r.db.Exec(ctx, queryUpdateTask,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	)

	if err != nil {
		return mapTaskDBError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

// Delete elimina una tarea por su ID
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, queryDeleteTask, id)
	if err != nil {
		return mapTaskDBError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

// ListByAssignee obtiene todas las tareas asignadas a un usuario
func (r *TaskRepository) ListByAssignee(ctx context.Context, userID string) ([]*domain.Task, error) {
	rows, err := r.db.Query(ctx, queryListTasksByAssignee, userID)
	if err != nil {
		return nil, mapTaskDBError(err)
	}
	defer rows.Close()

	tasks, err := scanTasks(rows)
	if err != nil {
		return nil, err
	}

	// Cargar comentarios para todas las tareas
	if err := loadTasksComments(ctx, r.db, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// ListAll obtiene todas las tareas ordenadas por fecha de creación
func (r *TaskRepository) ListAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.db.Query(ctx, queryListAllTasks)
	if err != nil {
		return nil, mapTaskDBError(err)
	}
	defer rows.Close()

	tasks, err := scanTasks(rows)
	if err != nil {
		return nil, err
	}

	// Cargar comentarios para todas las tareas
	if err := loadTasksComments(ctx, r.db, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddComment agrega un comentario a una tarea
func (r *TaskRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	_, err := r.db.Exec(ctx, queryInsertComment,
		comment.ID,
		comment.TaskID,
		comment.UserID,
		comment.Content,
		comment.CreatedAt,
	)

	return mapTaskDBError(err)
}
