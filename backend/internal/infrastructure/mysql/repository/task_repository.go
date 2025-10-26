package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
	"github.com/ryusuke/task_app_layerx/internal/infrastructure/mysql/model"
)

type taskRepository struct{}

// NewTaskRepositoryは新しいTaskRepository実装を作成する
func NewTaskRepository() domain.TaskRepository {
	return &taskRepository{}
}

// Createは新しいタスクをデータベースに挿入する
func (r *taskRepository) Create(ctx context.Context, ex domain.Executor, task *domain.Task) error {
	m := model.TaskFromDomain(task)

	query := `
		INSERT INTO tasks (owner_id, title, description, due_date, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := ex.ExecContext(ctx, query,
		m.OwnerID,
		m.Title,
		m.Description,
		m.DueDate,
		m.Status,
		m.Priority,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

// FindByIDはIDでタスクを取得する
func (r *taskRepository) FindByID(ctx context.Context, ex domain.Executor, taskID int64) (*domain.Task, error) {
	query := `
		SELECT id, owner_id, title, description, due_date, status, priority, created_at, updated_at, deleted_at
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL
	`

	row := ex.QueryRowContext(ctx, query, taskID)

	var m model.Task
	err := row.Scan(
		&m.ID,
		&m.OwnerID,
		&m.Title,
		&m.Description,
		&m.DueDate,
		&m.Status,
		&m.Priority,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	return m.ToDomain(), nil
}

// ListByUserIDはユーザーが所有または割り当てられているタスクを取得する
func (r *taskRepository) ListByUserID(ctx context.Context, ex domain.Executor, userID int64, limit, offset int) ([]*domain.Task, error) {
	// デフォルト値とバリデーション
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, owner_id, title, description, due_date, status, priority, created_at, updated_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL
		  AND (
		    owner_id = ?
		    OR EXISTS (
		      SELECT 1
		      FROM task_assignees
		      WHERE task_assignees.task_id = tasks.id
		        AND task_assignees.user_id = ?
		    )
		  )
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := ex.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var tasks []*domain.Task
	for rows.Next() {
		var m model.Task
		err := rows.Scan(
			&m.ID,
			&m.OwnerID,
			&m.Title,
			&m.Description,
			&m.DueDate,
			&m.Status,
			&m.Priority,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, m.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Updateは既存のタスクを更新する
func (r *taskRepository) Update(ctx context.Context, ex domain.Executor, task *domain.Task) error {
	m := model.TaskFromDomain(task)

	query := `
		UPDATE tasks
		SET title = ?, description = ?, due_date = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := ex.ExecContext(ctx, query,
		m.Title,
		m.Description,
		m.DueDate,
		m.Status,
		m.Priority,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

// Deleteはタスクの論理削除を実行する
func (r *taskRepository) Delete(ctx context.Context, ex domain.Executor, taskID int64, now time.Time) error {
	query := `
		UPDATE tasks
		SET deleted_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := ex.ExecContext(ctx, query, now, now, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
