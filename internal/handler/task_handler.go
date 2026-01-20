package handler

import (
	"encoding/json"
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/middleware"
	"gregperez/task-management-api/internal/service"
	"gregperez/task-management-api/pkg/helper"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req service.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), req, userID, userRole)
	if err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	task, err := h.taskService.GetTask(r.Context(), taskID, userID, userRole)
	if err != nil {
		helper.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req service.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	if err := h.taskService.UpdateTask(r.Context(), taskID, req, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "tarea actualizada"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	if err := h.taskService.DeleteTask(r.Context(), taskID, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "tarea eliminada"})
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req service.UpdateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	if err := h.taskService.UpdateTaskStatus(r.Context(), taskID, req, userID, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "estado actualizado"})
}

func (h *TaskHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	if err := h.taskService.AddComment(r.Context(), taskID, req.Content, userID, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusCreated, map[string]string{"message": "comentario agregado"})
}

func (h *TaskHandler) ListMyTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	tasks, err := h.taskService.ListUserTasks(r.Context(), userID)
	if err != nil {
		helper.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) ListAllTasks(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	tasks, err := h.taskService.ListAllTasks(r.Context(), userRole)
	if err != nil {
		helper.RespondError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, tasks)
}
