package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ryusuke/task_app_layerx/internal/domain"
	"github.com/ryusuke/task_app_layerx/pkg/auth"
	"github.com/ryusuke/task_app_layerx/pkg/hash"
)

type AuthUseCase struct {
	userRepo   domain.UserRepository
	txManager  domain.TxManager
	clock      domain.Clock
	jwtService auth.JWTService
	bcrypt     hash.BcryptService
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	txManager domain.TxManager,
	clock domain.Clock,
	jwtService auth.JWTService,
	bcrypt hash.BcryptService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		txManager:  txManager,
		clock:      clock,
		jwtService: jwtService,
		bcrypt:     bcrypt,
	}
}

func (u *AuthUseCase) Signup(ctx context.Context, req SignupRequest) (*AuthResponse, error) {
	// 入力の正規化
	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)

	// バリデーション
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return nil, domain.ErrPasswordTooShort
	}
	if name == "" {
		return nil, domain.ErrInvalidName
	}

	var response *AuthResponse

	err := u.txManager.Do(ctx, func(ctx context.Context, ex domain.Executor) error {
		// メールアドレスの重複チェック
		existingUser, err := u.userRepo.FindByEmail(ctx, ex, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}
		if existingUser != nil {
			return domain.ErrDuplicateEmail
		}

		// パスワードのハッシュ化
		hashedPassword, err := u.bcrypt.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// ユーザーエンティティ作成
		user, err := domain.NewUser(u.clock, email, name)
		if err != nil {
			return err
		}
		user.SetPasswordHash(u.clock, hashedPassword)

		// ユーザー保存
		if err := u.userRepo.Create(ctx, ex, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// JWTトークン生成
		token, err := u.jwtService.GenerateToken(user.ID, user.TokenVersion)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		response = &AuthResponse{
			Token: token,
			User: UserResponse{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// 入力の正規化
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// バリデーション
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}
	if req.Password == "" {
		return nil, domain.ErrInvalidPassword
	}

	// DBExecutorを作成（読み込み専用なのでトランザクション不要）
	executor := u.txManager.AsExecutor()

	// メールアドレスでユーザー検索
	user, err := u.userRepo.FindByEmail(ctx, executor, email)
	if err != nil {
		u.bcrypt.VerifyPassword("$2a$12$dummy", req.Password)
		return nil, domain.ErrUnauthorized
	}

	// パスワード検証
	if !u.bcrypt.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, domain.ErrUnauthorized
	}

	// JWTトークン生成
	token, err := u.jwtService.GenerateToken(user.ID, user.TokenVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, userID int64) error {
	return u.txManager.Do(ctx, func(ctx context.Context, ex domain.Executor) error {
		now := u.clock.Now()
		// TokenVersionをインクリメント → 既存のJWTが無効化される
		return u.userRepo.IncrementTokenVersion(ctx, ex, userID, now)
	})
}

func (u *AuthUseCase) GetUsers(ctx context.Context) ([]UserResponse, error) {
	executor := u.txManager.AsExecutor()
	users, err := u.userRepo.FindAll(ctx, executor)
	if err != nil {
		return nil, err
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}
	}

	return response, nil
}
