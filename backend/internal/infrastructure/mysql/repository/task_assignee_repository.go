package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ryusuke/task_app_layerx/internal/domain"
	"github.com/ryusuke/task_app_layerx/internal/infrastructure/mysql/model"
)

type taskAssigneeRepository struct{}

// NewTaskAssigneeRepository は新しい TaskAssigneeRepository 実装を作成します
func NewTaskAssigneeRepository() domain.TaskAssigneeRepository {
	return &taskAssigneeRepository{}
}

// Create は新しいタスク担当者をデータベースに挿入します
func (r *taskAssigneeRepository) Create(ctx context.Context, ex domain.Executor, assignee *domain.TaskAssignee) error {
	m := model.TaskAssigneeFromDomain(assignee)

	query := `
		INSERT INTO task_assignees (task_id, user_id, assigned_by, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := ex.ExecContext(ctx, query,
		m.TaskID,
		m.UserID,
		m.AssignedBy,
		m.CreatedAt,
	)
	if err != nil {
		// MySQL の Duplicate entry エラー（Error 1062）を検出
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "Error 1062") {
			return domain.ErrDuplicateAssignee
		}
		return fmt.Errorf("failed to create task assignee: %w", err)
	}

	return nil
}

// FindByTaskID は指定されたタスクのすべての担当者を取得します
func (r *taskAssigneeRepository) FindByTaskID(ctx context.Context, ex domain.Executor, taskID int64) ([]*domain.TaskAssignee, error) {
	query := `
		SELECT task_id, user_id, assigned_by, created_at
		FROM task_assignees
		WHERE task_id = ?
		ORDER BY created_at ASC
	`

	rows, err := ex.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task assignees: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var assignees []*domain.TaskAssignee
	for rows.Next() {
		var m model.TaskAssignee
		err := rows.Scan(
			&m.TaskID,
			&m.UserID,
			&m.AssignedBy,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task assignee: %w", err)
		}
		assignees = append(assignees, m.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task assignees: %w", err)
	}

	return assignees, nil
}

// DeleteByTaskID は指定されたタスクのすべての担当者を削除します
func (r *taskAssigneeRepository) DeleteByTaskID(ctx context.Context, ex domain.Executor, taskID int64) error {
	query := `
		DELETE FROM task_assignees
		WHERE task_id = ?
	`
	
	_, err := ex.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task assignees: %w", err)
	}

	return nil
}
