package domain

import "errors"

// 今回大きいPJや実務ではないのでerrorは一つのファイルにまとめる

// 共通エラー
var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
)

// User関連
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDuplicateEmail   = errors.New("email already exists")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrInvalidName      = errors.New("name is required")
	ErrNameTooLong      = errors.New("name must be less than 100 characters")
)

// Task関連
var (
	ErrTaskNotFound           = errors.New("task not found")
	ErrTitleRequired          = errors.New("title is required")
	ErrTitleTooLong           = errors.New("title must be less than 255 characters")
	ErrInvalidPriority        = errors.New("priority must be between 0 and 5")
	ErrInvalidStatus          = errors.New("invalid task status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// TaskAssignee関連
var (
	ErrDuplicateAssignee = errors.New("user already assigned to this task")
	ErrAssigneeNotFound  = errors.New("assignee not found")
)

// 認証関連
var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrTokenExpired = errors.New("token has expired")
)

// Validation補助構造体
type FieldError struct {
	Field   string
	Message string
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	return "validation error"
}