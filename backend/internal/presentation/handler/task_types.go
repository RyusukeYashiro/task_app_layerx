package handler

// CreateTaskRequestはタスク作成のリクエスト
type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"dueDate"`
	Priority    int     `json:"priority"`
	AssigneeIDs []int64 `json:"assigneeIds"`
}

// UpdateTaskRequestはタスク更新のリクエスト
type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"dueDate"`
	Status      *string `json:"status"`
	Priority    *int    `json:"priority"`
	AssigneeIDs []int64 `json:"assigneeIds"`
}

// TaskResponseはタスクのレスポンス
type TaskResponse struct {
	ID          int64              `json:"id"`
	OwnerID     int64              `json:"ownerId"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	DueDate     *string            `json:"dueDate"`
	Status      string             `json:"status"`
	Priority    int                `json:"priority"`
	Assignees   []AssigneeResponse `json:"assignees"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

// AssigneeResponseはアサイン情報のレスポンス
type AssigneeResponse struct {
	UserID     int64  `json:"userId"`
	AssignedBy int64  `json:"assignedBy"`
	AssignedAt string `json:"assignedAt"`
}

