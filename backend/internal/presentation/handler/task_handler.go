package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ryusuke/task_app_layerx/internal/presentation/middleware"
	taskuc "github.com/ryusuke/task_app_layerx/internal/usecase/task"
)

// TaskHandlerはタスク管理のHTTPハンドラー
type TaskHandler struct {
	taskUseCase *taskuc.TaskUseCase
}

// NewTaskHandlerで新しいTaskHandlerを作成
func NewTaskHandler(taskUseCase *taskuc.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// ListTasksはタスク一覧を取得
// GET /tasks
func (h *TaskHandler) ListTasks(c echo.Context) error {
	userID := middleware.GetUserID(c)

	// クエリパラメータを取得
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	req := taskuc.ListTasksRequest{
		Limit:  limit,
		Offset: offset,
	}

	resp, err := h.taskUseCase.ListTasks(c.Request().Context(), userID, req)
	if err != nil {
		return HandleError(c, err)
	}

	// レスポンス変換
	tasks := make([]TaskResponse, len(resp))
	for i, task := range resp {
		tasks[i] = toTaskResponse(task)
	}

	return c.JSON(http.StatusOK, tasks)
}

// CreateTaskはタスクを作成
// POST /tasks
func (h *TaskHandler) CreateTask(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		parsed, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "INVALID_DATE_FORMAT",
				Message: "dueDate must be in ISO8601 format",
			})
		}
		dueDate = &parsed
	}

	usecaseReq := taskuc.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		Priority:    req.Priority,
		AssigneeIDs: req.AssigneeIDs,
	}

	resp, err := h.taskUseCase.CreateTask(c.Request().Context(), userID, usecaseReq)
	if err != nil {
		return HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, toTaskResponse(resp))
}

// GetTaskはタスク詳細を取得
// GET /tasks/:id
func (h *TaskHandler) GetTask(c echo.Context) error {
	userID := middleware.GetUserID(c)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_TASK_ID",
			Message: "invalid task id",
		})
	}

	resp, err := h.taskUseCase.GetTask(c.Request().Context(), userID, taskID)
	if err != nil {
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, toTaskResponse(resp))
}

// UpdateTaskはタスクを更新
// PATCH /tasks/:id
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	userID := middleware.GetUserID(c)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_TASK_ID",
			Message: "invalid task id",
		})
	}

	var req UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
	}

	// DueDateをパース
	var dueDate *time.Time
	if req.DueDate != nil {
		parsed, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "INVALID_DATE_FORMAT",
				Message: "dueDate must be in ISO8601 format",
			})
		}
		dueDate = &parsed
	}

	usecaseReq := taskuc.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeIDs: req.AssigneeIDs,
	}

	resp, err := h.taskUseCase.UpdateTask(c.Request().Context(), userID, taskID, usecaseReq)
	if err != nil {
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, toTaskResponse(resp))
}

// DeleteTaskはタスクを削除
// DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	userID := middleware.GetUserID(c)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_TASK_ID",
			Message: "invalid task id",
		})
	}

	if err := h.taskUseCase.DeleteTask(c.Request().Context(), userID, taskID); err != nil {
		return HandleError(c, err)
	}

	// 204 No Content
	return c.NoContent(http.StatusNoContent)
}

// toTaskResponseはUseCaseのTaskResponseをHandlerのTaskResponseに変換
func toTaskResponse(task *taskuc.TaskResponse) TaskResponse {
	var dueDate *string
	if task.DueDate != nil {
		formatted := task.DueDate.Format(time.RFC3339)
		dueDate = &formatted
	}

	assignees := make([]AssigneeResponse, len(task.Assignees))
	for i, assignee := range task.Assignees {
		assignees[i] = AssigneeResponse{
			UserID:     assignee.UserID,
			AssignedBy: assignee.AssignedBy,
			AssignedAt: assignee.AssignedAt.Format(time.RFC3339),
		}
	}

	return TaskResponse{
		ID:          task.ID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     dueDate,
		Status:      task.Status,
		Priority:    task.Priority,
		Assignees:   assignees,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
}
