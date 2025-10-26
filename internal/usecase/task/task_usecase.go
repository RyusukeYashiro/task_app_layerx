package task

import (
	"context"
	"fmt"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// TaskUseCaseはタスク管理のユースケースを提供する
type TaskUseCase struct {
	taskRepo     domain.TaskRepository
	assigneeRepo domain.TaskAssigneeRepository
	userRepo     domain.UserRepository
	txManager    domain.TxManager
	clock        domain.Clock
}

// NewTaskUseCaseで新しいTaskUseCaseを作成
func NewTaskUseCase(
	taskRepo domain.TaskRepository,
	assigneeRepo domain.TaskAssigneeRepository,
	userRepo domain.UserRepository,
	txManager domain.TxManager,
	clock domain.Clock,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:     taskRepo,
		assigneeRepo: assigneeRepo,
		userRepo:     userRepo,
		txManager:    txManager,
		clock:        clock,
	}
}

// ListTasksはユーザーに関連するタスク一覧を取得
func (u *TaskUseCase) ListTasks(ctx context.Context, userID int64, req ListTasksRequest) ([]*TaskResponse, error) {
	// デフォルト値を設定
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	executor := u.txManager.AsExecutor()

	// ユーザーに関連するタスク一覧を取得
	tasks, err := u.taskRepo.ListByUserID(ctx, executor, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// レスポンスを作成
	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		// タスクに関連するアサイン一覧を取得
		assignees, err := u.assigneeRepo.FindByTaskID(ctx, executor, task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find assignees: %w", err)
		}

		responses[i] = &TaskResponse{
			ID:          task.ID,
			OwnerID:     task.OwnerID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      string(task.Status),
			Priority:    task.Priority,
			Assignees:   toAssigneeResponses(assignees),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return responses, nil
}

// CreateTaskはタスクを作成
func (u *TaskUseCase) CreateTask(ctx context.Context, userID int64, req CreateTaskRequest) (*TaskResponse, error) {
	var response *TaskResponse

	err := u.txManager.Do(ctx, func(ctx context.Context, ex domain.Executor) error {
		// タスクエンティティを作成
		task, err := domain.NewTask(u.clock, userID, req.Title)
		if err != nil {
			return err
		}

		// オプション項目を設定
		if req.Description != nil {
			task.UpdateDescription(u.clock, req.Description)
		}
		if req.DueDate != nil {
			task.UpdateDueDate(u.clock, req.DueDate)
		}
		if req.Priority != 0 {
			task.Priority = req.Priority
			if err := task.ValidatePriority(); err != nil {
				return err
			}
		}

		// タスクを保存
		if err := u.taskRepo.Create(ctx, ex, task); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		// アサイン処理
		assignees := make([]*domain.TaskAssignee, 0)
		for _, assigneeID := range req.AssigneeIDs {
			// アサイン先ユーザーが存在するか確認
			_, err := u.userRepo.FindByID(ctx, ex, assigneeID)
			if err != nil {
				return fmt.Errorf("assignee user not found: %w", err)
			}

			assignee := &domain.TaskAssignee{
				TaskID:     task.ID,
				UserID:     assigneeID,
				AssignedBy: userID,
				CreatedAt:  u.clock.Now(),
			}

			if err := u.assigneeRepo.Create(ctx, ex, assignee); err != nil {
				return fmt.Errorf("failed to create assignee: %w", err)
			}

			assignees = append(assignees, assignee)
		}

		response = &TaskResponse{
			ID:          task.ID,
			OwnerID:     task.OwnerID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      string(task.Status),
			Priority:    task.Priority,
			Assignees:   toAssigneeResponses(assignees),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetTask はタスク詳細を取得する
func (u *TaskUseCase) GetTask(ctx context.Context, userID, taskID int64) (*TaskResponse, error) {
	executor := u.txManager.AsExecutor()

	// タスク取得
	task, err := u.taskRepo.FindByID(ctx, executor, taskID)
	if err != nil {
		return nil, err
	}

	// 権限チェック（オーナーまたはアサイン先のみ）
	if !task.IsOwner(userID) {
		// アサイン先かチェック
		assignees, err := u.assigneeRepo.FindByTaskID(ctx, executor, taskID)
		if err != nil {
			return nil, fmt.Errorf("failed to find assignees: %w", err)
		}

		isAssignee := false
		for _, assignee := range assignees {
			if assignee.UserID == userID {
				isAssignee = true
				break
			}
		}

		if !isAssignee {
			// オーナーでもアサイン先でもない場合は存在を隠蔽
			return nil, domain.ErrTaskNotFound
		}
	}

	// アサイン一覧を取得
	assignees, err := u.assigneeRepo.FindByTaskID(ctx, executor, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find assignees: %w", err)
	}

	return &TaskResponse{
		ID:          task.ID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      string(task.Status),
		Priority:    task.Priority,
		Assignees:   toAssigneeResponses(assignees),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

// UpdateTaskはタスクを更新
func (u *TaskUseCase) UpdateTask(ctx context.Context, userID, taskID int64, req UpdateTaskRequest) (*TaskResponse, error) {
	var response *TaskResponse

	err := u.txManager.Do(ctx, func(ctx context.Context, ex domain.Executor) error {
		// タスクを取得
		task, err := u.taskRepo.FindByID(ctx, ex, taskID)
		if err != nil {
			return err
		}

		// 権限チェック（オーナーのみ更新可能）
		if !task.IsOwner(userID) {
			return domain.ErrForbidden
		}

		// 各項目を更新
		if req.Title != nil {
			if err := task.UpdateTitle(u.clock, *req.Title); err != nil {
				return err
			}
		}

		if req.Description != nil {
			task.UpdateDescription(u.clock, req.Description)
		}

		if req.DueDate != nil {
			task.UpdateDueDate(u.clock, req.DueDate)
		}

		if req.Status != nil {
			newStatus := domain.TaskStatus(*req.Status)
			if err := task.ValidateStatusTransaction(newStatus); err != nil {
				return err
			}
			task.Status = newStatus
			task.UpdatedAt = u.clock.Now()
		}

		if req.Priority != nil {
			if err := task.UpdatePriority(u.clock, *req.Priority); err != nil {
				return err
			}
		}

		// タスクを保存
		if err := u.taskRepo.Update(ctx, ex, task); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		// アサインを更新（指定されている場合）
		var assignees []*domain.TaskAssignee
		if req.AssigneeIDs != nil {
			// 既存のアサインを削除（idempotentな操作）
			if err := u.assigneeRepo.DeleteByTaskID(ctx, ex, taskID); err != nil {
				return fmt.Errorf("failed to delete assignees: %w", err)
			}

			// 新しいアサインを作成
			for _, assigneeID := range req.AssigneeIDs {
				// アサイン先ユーザーが存在するか確認
				_, err := u.userRepo.FindByID(ctx, ex, assigneeID)
				if err != nil {
					return fmt.Errorf("assignee user not found: %w", err)
				}

				assignee := &domain.TaskAssignee{
					TaskID:     taskID,
					UserID:     assigneeID,
					AssignedBy: userID,
					CreatedAt:  u.clock.Now(),
				}

				if err := u.assigneeRepo.Create(ctx, ex, assignee); err != nil {
					return fmt.Errorf("failed to create assignee: %w", err)
				}

				assignees = append(assignees, assignee)
			}
		} else {
			// AssigneeIDsが指定されていない場合は既存のアサインを維持
			assignees, err = u.assigneeRepo.FindByTaskID(ctx, ex, taskID)
			if err != nil {
				return fmt.Errorf("failed to find assignees: %w", err)
			}
		}

		// レスポンスを作成
		response = &TaskResponse{
			ID:          task.ID,
			OwnerID:     task.OwnerID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      string(task.Status),
			Priority:    task.Priority,
			Assignees:   toAssigneeResponses(assignees),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteTaskはタスクを削除
func (u *TaskUseCase) DeleteTask(ctx context.Context, userID, taskID int64) error {
	return u.txManager.Do(ctx, func(ctx context.Context, ex domain.Executor) error {
		// タスクを取得
		task, err := u.taskRepo.FindByID(ctx, ex, taskID)
		if err != nil {
			return err
		}

		// 権限チェック（オーナーのみ削除可能）
		if !task.IsOwner(userID) {
			return domain.ErrForbidden
		}

		// アサインを削除
		if err := u.assigneeRepo.DeleteByTaskID(ctx, ex, taskID); err != nil {
			return fmt.Errorf("failed to delete assignees: %w", err)
		}

		// タスクを削除
		now := u.clock.Now()
		if err := u.taskRepo.Delete(ctx, ex, taskID, now); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		return nil
	})
}

// toAssigneeResponsesはdomain.TaskAssigneeのスライスをAssigneeResponseのスライスに変換
func toAssigneeResponses(assignees []*domain.TaskAssignee) []AssigneeResponse {
	responses := make([]AssigneeResponse, len(assignees))
	for i, assignee := range assignees {
		responses[i] = AssigneeResponse{
			UserID:     assignee.UserID,
			AssignedBy: assignee.AssignedBy,
			AssignedAt: assignee.CreatedAt,
		}
	}
	return responses
}
