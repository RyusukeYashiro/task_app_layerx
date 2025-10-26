package domain_test

import (
	"testing"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

func TestNewTask(t *testing.T) {
	clock := &mockClock{}

	tests := []struct {
		name      string
		ownerID   int64
		title     string
		wantError bool
	}{
		{
			name:      "正常なタスク作成",
			ownerID:   1,
			title:     "テストタスク",
			wantError: false,
		},
		{
			name:      "タイトルが空",
			ownerID:   1,
			title:     "",
			wantError: true,
		},
		{
			name:      "タイトルが長すぎる",
			ownerID:   1,
			title:     string(make([]byte, 256)),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := domain.NewTask(clock, tt.ownerID, tt.title)
			if tt.wantError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが返されませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが期待されていませんでした: %v", err)
				}
				if task == nil {
					t.Error("タスクがnilです")
				}
				if task.Status != domain.TaskStatusTODO {
					t.Errorf("Status = %v, want %v", task.Status, domain.TaskStatusTODO)
				}
			}
		})
	}
}

func TestTask_UpdateTitle(t *testing.T) {
	clock := &mockClock{}
	task, _ := domain.NewTask(clock, 1, "元のタイトル")

	tests := []struct {
		name      string
		newTitle  string
		wantError bool
	}{
		{
			name:      "正常なタイトル更新",
			newTitle:  "新しいタイトル",
			wantError: false,
		},
		{
			name:      "空のタイトル",
			newTitle:  "",
			wantError: true,
		},
		{
			name:      "長すぎるタイトル",
			newTitle:  string(make([]byte, 256)),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.UpdateTitle(clock, tt.newTitle)
			if tt.wantError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが返されませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが期待されていませんでした: %v", err)
				}
			}
		})
	}
}

func TestTask_UpdatePriority(t *testing.T) {
	clock := &mockClock{}
	task, _ := domain.NewTask(clock, 1, "テスト")

	tests := []struct {
		name      string
		priority  int
		wantError bool
	}{
		{name: "優先度0", priority: 0, wantError: false},
		{name: "優先度3", priority: 3, wantError: false},
		{name: "優先度5", priority: 5, wantError: false},
		{name: "優先度が負", priority: -1, wantError: true},
		{name: "優先度が大きすぎる", priority: 6, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.UpdatePriority(clock, tt.priority)
			if tt.wantError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが返されませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが期待されていませんでした: %v", err)
				}
				if task.Priority != tt.priority {
					t.Errorf("Priority = %v, want %v", task.Priority, tt.priority)
				}
			}
		})
	}
}

func TestTask_IsOwner(t *testing.T) {
	clock := &mockClock{}
	task, _ := domain.NewTask(clock, 1, "テスト")

	tests := []struct {
		name   string
		userID int64
		want   bool
	}{
		{name: "オーナー", userID: 1, want: true},
		{name: "オーナーではない", userID: 2, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := task.IsOwner(tt.userID)
			if got != tt.want {
				t.Errorf("IsOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}
