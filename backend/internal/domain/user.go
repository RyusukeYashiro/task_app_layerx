package domain

import (
	"regexp"
	"strings"
	"time"
)

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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewUser新しいユーザーを作成
func NewUser(clock Clock, email, name string) (*User, error) {
	now := clock.Now()

	u := &User{
		Email:        normalizeEmail(email),
		Name:         strings.TrimSpace(name),
		TokenVersion: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.ValidateEmail(); err != nil {
		return nil, err
	}
	if err := u.ValidateName(); err != nil {
		return nil, err
	}
	return u, nil
}

// ValidateEmailメールアドレスを検証
func (u *User) ValidateEmail() error {
	if u.Email == "" || !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateNameは名前を検証
func (u *User) ValidateName() error {
	if strings.TrimSpace(u.Name) == "" {
		return ErrUserNotFound
	}
	if len(u.Name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// IncrementTokenVersionはトークンバージョンをインクリメント
func (u *User) IncrementTokenVersion(clock Clock) {
	u.TokenVersion++
	u.UpdatedAt = clock.Now()
}

// IsDeletedは削除フラグを確認
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SetPasswordHashはパスワードハッシュを設定
func (u *User) SetPasswordHash(clock Clock, hash string) {
	u.PasswordHash = hash
	u.UpdatedAt = clock.Now()
}

// normalizeEmailはメールアドレスを正規化
func normalizeEmail(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}