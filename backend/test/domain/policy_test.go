package domain_test

import (
	"testing"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

func TestCanViewTask(t *testing.T) {
	clock := &mockClock{}
	task, _ := domain.NewTask(clock, 1, "test task")

	tests := []struct {
		name      string
		assignees []*domain.TaskAssignee
		userID    int64
		want      bool
	}{
		{
			name:      "オーナーは閲覧可能",
			assignees: []*domain.TaskAssignee{},
			userID:    1,
			want:      true,
		},
		{
			name: "アサイン先は閲覧可能",
			assignees: []*domain.TaskAssignee{
				{UserID: 2},
			},
			userID: 2,
			want:   true,
		},
		{
			name:      "オーナーでもアサイン先でもない",
			assignees: []*domain.TaskAssignee{},
			userID:    999,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.CanViewTask(task, tt.assignees, tt.userID)
			if got != tt.want {
				t.Errorf("CanViewTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanEditTask(t *testing.T) {
	clock := &mockClock{}
	task, _ := domain.NewTask(clock, 1, "test task")

	tests := []struct {
		name   string
		userID int64
		want   bool
	}{
		{name: "オーナーは編集可能", userID: 1, want: true},
		{name: "オーナー以外は編集不可", userID: 2, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.CanEditTask(task, tt.userID)
			if got != tt.want {
				t.Errorf("CanEditTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
