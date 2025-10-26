package model

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// Userはusersテーブルの構造を表す
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

// ToDomainはDBモデルをドメインエンティティに変換
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

// UserFromDomainはドメインエンティティをDBモデルに変換
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
