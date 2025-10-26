package auth_test

import (
	"testing"
	"time"

	"github.com/ryusuke/task_app_layerx/pkg/auth"
)

func TestJWTService_GenerateToken(t *testing.T) {
	secret := "test-secret-key"
	issuer := "test-issuer"
	expiration := 1 * time.Hour
	now := time.Now()

	service := auth.NewJWTService(secret, issuer, expiration, func() time.Time {
		return now
	})

	userID := int64(123)
	tokenVersion := 1

	token, err := service.GenerateToken(userID, tokenVersion)
	if err != nil {
		t.Fatalf("トークン生成に失敗しました: %v", err)
	}

	if token == "" {
		t.Error("トークンが空です")
	}
}

func TestJWTService_ParseToken(t *testing.T) {
	secret := "test-secret-key"
	issuer := "test-issuer"
	expiration := 1 * time.Hour
	now := time.Now()

	service := auth.NewJWTService(secret, issuer, expiration, func() time.Time {
		return now
	})

	userID := int64(123)
	tokenVersion := 1

	token, _ := service.GenerateToken(userID, tokenVersion)

	t.Run("有効なトークン", func(t *testing.T) {
		claims, err := service.ParseToken(token)
		if err != nil {
			t.Fatalf("トークン検証に失敗しました: %v", err)
		}

		if claims.UID != userID {
			t.Errorf("UID = %v, want %v", claims.UID, userID)
		}

		if claims.TokenVersion != tokenVersion {
			t.Errorf("TokenVersion = %v, want %v", claims.TokenVersion, tokenVersion)
		}
	})

	t.Run("無効なトークン", func(t *testing.T) {
		_, err := service.ParseToken("invalid.token.here")
		if err == nil {
			t.Error("無効なトークンでエラーが期待されましたが、成功しました")
		}
	})

	t.Run("異なるシークレット", func(t *testing.T) {
		otherService := auth.NewJWTService("different-secret", issuer, expiration, func() time.Time {
			return now
		})
		otherToken, _ := otherService.GenerateToken(userID, tokenVersion)

		_, err := service.ParseToken(otherToken)
		if err == nil {
			t.Error("異なるシークレットで生成されたトークンでエラーが期待されました")
		}
	})
}
