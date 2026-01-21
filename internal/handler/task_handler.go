package handler

import (
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/handler/request"
	"gregperez/task-management-api/internal/service"
	"gregperez/task-management-api/internal/service/dto"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService service.TaskServiceInterface
}

func NewTaskHandler(taskService service.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateTask crea una nueva tarea (solo Admin)
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, userRole, err := GetUserContext(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.CreateTaskRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := dto.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		AssignedTo:  req.AssignedTo,
	}

	task, err := h.taskService.CreateTask(ctx, serviceReq, userID, userRole)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusCreated, task)
}

// GetTask obtiene una tarea por ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	userID, userRole, err := GetUserContext(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	task, err := h.taskService.GetTask(ctx, taskID, userID, userRole)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, task)
}

// UpdateTask actualiza una tarea existente (solo Admin)
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.UpdateTaskRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := dto.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	if err := h.taskService.UpdateTask(ctx, taskID, serviceReq, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "tarea actualizada")
}

// DeleteTask elimina una tarea (solo Admin)
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	if err := h.taskService.DeleteTask(ctx, taskID, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "tarea eliminada")
}

// UpdateTaskStatus actualiza el estado de una tarea (Ejecutor)
func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	userID, userRole, err := GetUserContext(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.UpdateTaskStatusRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := dto.UpdateTaskStatusRequest{
		Status: domain.TaskStatus(req.Status),
	}

	if err := h.taskService.UpdateTaskStatus(ctx, taskID, serviceReq, userID, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "estado actualizado")
}

// AddComment agrega un comentario a una tarea (Ejecutor)
func (h *TaskHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	userID, userRole, err := GetUserContext(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.AddCommentRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	if err := h.taskService.AddComment(ctx, taskID, req.Content, userID, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusCreated, "comentario agregado")
}

// ListMyTasks lista las tareas del usuario autenticado (Ejecutor)
func (h *TaskHandler) ListMyTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := GetUserID(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	tasks, err := h.taskService.ListUserTasks(ctx, userID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, tasks)
}

// ListAllTasks lista todas las tareas (Auditor)
func (h *TaskHandler) ListAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	tasks, err := h.taskService.ListAllTasks(ctx, userRole)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, tasks)
}
