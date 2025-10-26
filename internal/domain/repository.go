package domain

import "context"

// UserRepository defines operations for User persistence
type UserRepository interface {
	Create(ctx context.Context, ex Executor, user *User) error
	FindByID(ctx context.Context, ex Executor, id int64) (*User, error)
	FindByEmail(ctx context.Context, ex Executor, email string) (*User, error)
	Update(ctx context.Context, ex Executor, user *User) error
	IncrementTokenVersion(ctx context.Context, ex Executor, userID int64) error
}

// TaskRepository defines operations for Task persistence
type TaskRepository interface {
	Create(ctx context.Context, ex Executor, task *Task) error
	FindByID(ctx context.Context, ex Executor, taskID int64) (*Task, error)
	ListByUserID(ctx context.Context, ex Executor, userID int64) ([]*Task, error)
	Update(ctx context.Context, ex Executor, task *Task) error
	Delete(ctx context.Context, ex Executor, taskID int64) error
}

// TaskAssigneeRepository defines operations for TaskAssignee persistence
type TaskAssigneeRepository interface {
	Create(ctx context.Context, ex Executor, assignee *TaskAssignee) error
	FindByTaskID(ctx context.Context, ex Executor, taskID int64) ([]*TaskAssignee, error)
	DeleteByTaskID(ctx context.Context, ex Executor, taskID int64) error
}
