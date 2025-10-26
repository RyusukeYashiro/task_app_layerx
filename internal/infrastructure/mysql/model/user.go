package model

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// User は users テーブルの構造を表します
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	TokenVersion int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// ToDomain はDBモデルをドメインエンティティに変換します
func (m *User) ToDomain() *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Name:         m.Name,
		TokenVersion: m.TokenVersion,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
}

// UserFromDomain はドメインエンティティをDBモデルに変換します
func UserFromDomain(u *domain.User) *User {
	return &User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		TokenVersion: u.TokenVersion,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
