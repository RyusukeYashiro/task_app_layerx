package domain

import (
	"regexp"
	"strings"
	"time"
)

// ドメイン内で使う時計（テスト容易性のための注入点）
type Clock interface {
	Now() time.Time
}

type User struct {
	ID           uint64
	Email        string
	PasswordHash string
	Name         string
	TokenVersion int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewUser 新しいユーザーを作成
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

// ValidateEmail メールアドレスを検証
func (u *User) ValidateEmail() error {
	if u.Email == "" || !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateName 名前を検証
func (u *User) ValidateName() error {
	if strings.TrimSpace(u.Name) == "" {
		return ErrUserNotFound
	}
	if len(u.Name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// IncrementTokenVersion トークンバージョンをインクリメント
func (u *User) IncrementTokenVersion(clock Clock) {
	u.TokenVersion++
	u.UpdatedAt = clock.Now()
}

// IsDeleted 削除フラグを確認
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SetPasswordHash パスワードハッシュを設定
func (u *User) SetPasswordHash(clock Clock, hash string) {
	u.PasswordHash = hash
	u.UpdatedAt = clock.Now()
}

// normalizeEmail メールアドレスを正規化
func normalizeEmail(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}