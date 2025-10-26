package domain_test

import (
	"testing"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

func TestNewTaskAssignee(t *testing.T) {
	clock := &mockClock{}

	tests := []struct {
		name       string
		taskID     int64
		userID     int64
		assignedBy int64
		wantErr    bool
	}{
		{
			name:       "正常なタスク担当者作成",
			taskID:     1,
			userID:     2,
			assignedBy: 1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignee, err := domain.NewTaskAssignee(clock, tt.taskID, tt.userID, tt.assignedBy)
			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが返されませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが期待されていませんでした: %v", err)
				}
				if assignee == nil {
					t.Error("担当者がnilです")
				}
				if assignee.TaskID != tt.taskID {
					t.Errorf("TaskID = %v, want %v", assignee.TaskID, tt.taskID)
				}
				if assignee.UserID != tt.userID {
					t.Errorf("UserID = %v, want %v", assignee.UserID, tt.userID)
				}
				if assignee.AssignedBy != tt.assignedBy {
					t.Errorf("AssignedBy = %v, want %v", assignee.AssignedBy, tt.assignedBy)
				}
			}
		})
	}
}
