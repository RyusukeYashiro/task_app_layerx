package model

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// TaskAssigneeはtask_assigneesテーブルの構造を現す
type TaskAssignee struct {
	TaskID     int64
	UserID     int64
	AssignedBy int64
	CreatedAt  time.Time
}

// ToDomainはDBモデルをドメインエンティティに変換
func (m *TaskAssignee) ToDomain() *domain.TaskAssignee {
	return &domain.TaskAssignee{
		TaskID:     m.TaskID,
		UserID:     m.UserID,
		AssignedBy: m.AssignedBy,
		CreatedAt:  m.CreatedAt,
	}
}

// TaskAssigneeFromDomainはドメインエンティティをDBモデルに変換
func TaskAssigneeFromDomain(ta *domain.TaskAssignee) *TaskAssignee {
	return &TaskAssignee{
		TaskID:     ta.TaskID,
		UserID:     ta.UserID,
		AssignedBy: ta.AssignedBy,
		CreatedAt:  ta.CreatedAt,
	}
}
