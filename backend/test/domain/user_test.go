package domain_test

import (
	"testing"
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

type mockClock struct {
	now time.Time
}

func (m *mockClock) Now() time.Time {
	return m.now
}

func TestNewUser(t *testing.T) {
	clock := &mockClock{now: time.Now()}

	tests := []struct {
		name      string
		email     string
		userName  string
		wantError bool
	}{
		{
			name:      "正常なユーザー作成",
			email:     "test@example.com",
			userName:  "テストユーザー",
			wantError: false,
		},
		{
			name:      "メールアドレスが空",
			email:     "",
			userName:  "テストユーザー",
			wantError: true,
		},
		{
			name:      "名前が空",
			email:     "test@example.com",
			userName:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(clock, tt.email, tt.userName)
			if tt.wantError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが返されませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが期待されていませんでした: %v", err)
				}
				if user == nil {
					t.Error("ユーザーがnilです")
				}
			}
		})
	}
}
