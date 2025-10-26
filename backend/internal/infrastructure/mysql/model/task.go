package model

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// Taskはtasksテーブルの構造を現す
type Task struct {
	ID          int64
	OwnerID     int64
	Title       string
	Description *string
	DueDate     *time.Time
	Status      string
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ToDomainはDBモデルをドメインエンティティに変換
func (m *Task) ToDomain() *domain.Task {
	// ステータス値の安全化：未知の値はデフォルト（TODO）にフォールバック
	status := domain.TaskStatus(m.Status)
	if !isValidStatus(status) {
		status = domain.TaskStatusTODO
	}

	return &domain.Task{
		ID:          m.ID,
		OwnerID:     m.OwnerID,
		Title:       m.Title,
		Description: m.Description,
		DueDate:     m.DueDate,
		Status:      status,
		Priority:    m.Priority,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

// isValidStatusはステータス値が有効かどうかをチェック
func isValidStatus(status domain.TaskStatus) bool {
	return status == domain.TaskStatusTODO ||
		status == domain.TaskStatusIN_PROGRESS ||
		status == domain.TaskStatusDONE
}

// TaskFromDomainはドメインエンティティをDBモデルに変換
func TaskFromDomain(t *domain.Task) *Task {
	return &Task{
		ID:          t.ID,
		OwnerID:     t.OwnerID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Status:      string(t.Status),
		Priority:    t.Priority,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
	}
}
