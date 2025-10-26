package hash_test

import (
	"testing"

	"github.com/ryusuke/task_app_layerx/pkg/hash"
)

func TestBcryptService(t *testing.T) {
	service := hash.NewBcryptService(10)
	password := "testpassword123"

	t.Run("パスワードのハッシュ化", func(t *testing.T) {
		hashed, err := service.HashPassword(password)
		if err != nil {
			t.Fatalf("パスワードのハッシュ化に失敗しました: %v", err)
		}

		if hashed == "" {
			t.Error("ハッシュが空です")
		}

		if hashed == password {
			t.Error("ハッシュが平文と同じです")
		}
	})

	t.Run("正しいパスワードの検証", func(t *testing.T) {
		hashed, _ := service.HashPassword(password)
		
		if !service.VerifyPassword(hashed, password) {
			t.Error("正しいパスワードの検証に失敗しました")
		}
	})

	t.Run("間違ったパスワードの検証", func(t *testing.T) {
		hashed, _ := service.HashPassword(password)
		
		if service.VerifyPassword(hashed, "wrongpassword") {
			t.Error("間違ったパスワードが検証されてしまいました")
		}
	})
}
