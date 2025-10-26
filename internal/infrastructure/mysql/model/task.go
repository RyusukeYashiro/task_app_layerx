package model

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// Task は tasks テーブルの構造を表します
type Task struct {
	ID          int64
	OwnerID     int64
	Title       string
	Description *string
	DueDate     *time.Time
	Status      string // DB ENUM -> string
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ToDomain はDBモデルをドメインエンティティに変換します
func (m *Task) ToDomain() *domain.Task {
	return &domain.Task{
		ID:          m.ID,
		OwnerID:     m.OwnerID,
		Title:       m.Title,
		Description: m.Description,
		DueDate:     m.DueDate,
		Status:      domain.TaskStatus(m.Status),
		Priority:    m.Priority,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

// TaskFromDomain はドメインエンティティをDBモデルに変換します
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
