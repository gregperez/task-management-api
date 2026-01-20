package postgres

// Constantes de columnas para evitar duplicación
const (
	taskColumns    = "id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at"
	commentColumns = "id, task_id, user_id, content, created_at"
)

// Queries SQL para operaciones de tareas
const (
	queryInsertTask = `
		INSERT INTO tasks (id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	querySelectTaskByID = `
		SELECT ` + taskColumns + `
		FROM tasks
		WHERE id = $1
	`

	queryUpdateTask = `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, updated_at = $5
		WHERE id = $6
	`

	queryDeleteTask = `
		DELETE FROM tasks WHERE id = $1
	`

	queryListTasksByAssignee = `
		SELECT ` + taskColumns + `
		FROM tasks
		WHERE assigned_to = $1
		ORDER BY due_date ASC
	`

	queryListAllTasks = `
		SELECT ` + taskColumns + `
		FROM tasks
		ORDER BY created_at DESC
	`
)

// Queries SQL para operaciones de comentarios
const (
	queryInsertComment = `
		INSERT INTO comments (id, task_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	querySelectCommentsByTaskID = `
		SELECT ` + commentColumns + `
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC
	`
)
