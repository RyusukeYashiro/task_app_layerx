package task

import "time"

// ListTasksRequest はタスク一覧取得のリクエスト
type ListTasksRequest struct {
	Limit  int
	Offset int
}

// CreateTaskRequest はタスク作成のリクエスト
type CreateTaskRequest struct {
	Title       string
	Description *string
	DueDate     *time.Time
	Priority    int
	AssigneeIDs []int64
}

// UpdateTaskRequest はタスク更新のリクエスト
type UpdateTaskRequest struct {
	Title       *string
	Description *string
	DueDate     *time.Time
	Status      *string
	Priority    *int
	AssigneeIDs []int64
}

// TaskResponse はタスクのレスポンス
type TaskResponse struct {
	ID          int64
	OwnerID     int64
	Title       string
	Description *string
	DueDate     *time.Time
	Status      string
	Priority    int
	Assignees   []AssigneeResponse
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AssigneeResponse はアサイン情報のレスポンス
type AssigneeResponse struct {
	UserID     int64
	AssignedBy int64
	AssignedAt time.Time
}
