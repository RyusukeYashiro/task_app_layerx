package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
	"github.com/ryusuke/task_app_layerx/internal/infrastructure/mysql/model"
)

type userRepository struct{}

// NewUserRepositoryは新しいUserRepository実装を作成する
func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}

// Createは新しいユーザーをデータベースに挿入する
func (r *userRepository) Create(ctx context.Context, ex domain.Executor, user *domain.User) error {
	m := model.UserFromDomain(user)

	query := `
		INSERT INTO users (email, password_hash, name, token_version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := ex.ExecContext(ctx, query,
		m.Email,
		m.PasswordHash,
		m.Name,
		m.TokenVersion,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

// FindByIDはIDでユーザーを取得する
func (r *userRepository) FindByID(ctx context.Context, ex domain.Executor, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, token_version, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	row := ex.QueryRowContext(ctx, query, id)

	var m model.User
	err := row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.Name,
		&m.TokenVersion,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return m.ToDomain(), nil
}

// FindByEmailはメールアドレスでユーザーを取得する
func (r *userRepository) FindByEmail(ctx context.Context, ex domain.Executor, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, token_version, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	row := ex.QueryRowContext(ctx, query, email)

	var m model.User
	err := row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.Name,
		&m.TokenVersion,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return m.ToDomain(), nil
}

// Updateは既存のユーザーを更新する
func (r *userRepository) Update(ctx context.Context, ex domain.Executor, user *domain.User) error {
	m := model.UserFromDomain(user)

	query := `
		UPDATE users
		SET email = ?, password_hash = ?, name = ?, token_version = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := ex.ExecContext(ctx, query,
		m.Email,
		m.PasswordHash,
		m.Name,
		m.TokenVersion,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// IncrementTokenVersionはJWT無効化のためにトークンバージョンをインクリメントする
func (r *userRepository) IncrementTokenVersion(ctx context.Context, ex domain.Executor, userID int64, updatedAt time.Time) error {
	query := `
		UPDATE users
		SET token_version = token_version + 1, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := ex.ExecContext(ctx, query, updatedAt, userID)
	if err != nil {
		return fmt.Errorf("failed to increment token version: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
