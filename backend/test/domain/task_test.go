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
