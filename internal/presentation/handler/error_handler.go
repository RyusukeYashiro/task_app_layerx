package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// ErrorResponseはエラーレスポンスの構造
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HandleErrorはDomainエラーをHTTPエラーに変換
func HandleError(c echo.Context, err error) error {
	if errors.Is(err, domain.ErrUnauthorized) {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}
	// トークンエラー (401)
	if errors.Is(err, domain.ErrInvalidToken) {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: "invalid or expired token",
		})
	}
	// トークンが期限切れ (401)	
	if errors.Is(err, domain.ErrTokenExpired) {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "TOKEN_EXPIRED",
			Message: "token has expired",
		})
	}

	// 権限がない (403)
	if errors.Is(err, domain.ErrForbidden) {
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FORBIDDEN",
			Message: "forbidden",
		})
	}

	// リソースが見つからない (404)
	if errors.Is(err, domain.ErrNotFound) ||
		errors.Is(err, domain.ErrUserNotFound) ||
		errors.Is(err, domain.ErrTaskNotFound) {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		})
	}

	// メールアドレスが無効 (400)
	if errors.Is(err, domain.ErrInvalidEmail) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "invalid email format",
			Details: map[string]interface{}{"field": "email"},
		})
	}
	// パスワードが無効 (400)
	if errors.Is(err, domain.ErrInvalidPassword) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "invalid password",
			Details: map[string]interface{}{"field": "password"},
		})
	}
	// パスワードが短すぎる (400)
	if errors.Is(err, domain.ErrPasswordTooShort) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "password must be at least 8 characters",
			Details: map[string]interface{}{"field": "password"},
		})
	}
	// 名前が無効 (400)
	if errors.Is(err, domain.ErrInvalidName) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "name is required",
			Details: map[string]interface{}{"field": "name"},
		})
	}
	// タイトルが無効 (400)
	if errors.Is(err, domain.ErrTitleRequired) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "title is required",
			Details: map[string]interface{}{"field": "title"},
		})
	}
	// タイトルが長すぎる (400)
	if errors.Is(err, domain.ErrTitleTooLong) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "title must be less than 255 characters",
			Details: map[string]interface{}{"field": "title"},
		})
	}
	// 優先度が無効 (400)
	if errors.Is(err, domain.ErrInvalidPriority) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "priority must be between 0 and 5",
			Details: map[string]interface{}{"field": "priority"},
		})
	}
	// ステータス遷移が無効 (400)
	if errors.Is(err, domain.ErrInvalidStatusTransition) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "invalid status transition",
			Details: map[string]interface{}{"field": "status"},
		})
	}

	// メールアドレスが重複 (409)
	if errors.Is(err, domain.ErrDuplicateEmail) {
		return c.JSON(http.StatusConflict, ErrorResponse{
			Code:    "CONFLICT",
			Message: "email already exists",
			Details: map[string]interface{}{"field": "email"},
		})
	}
	// アサイン先が重複 (409)
	if errors.Is(err, domain.ErrDuplicateAssignee) {
		return c.JSON(http.StatusConflict, ErrorResponse{
			Code:    "CONFLICT",
			Message: "user already assigned to this task",
		})
	}

	// 内部エラー (500)
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
	})
}
