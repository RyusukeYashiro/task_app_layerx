package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

type txManager struct {
	db *sql.DB
}

// NewTxManagerは新しいTxManager実装を作成する
func NewTxManager(db *sql.DB) domain.TxManager {
	return &txManager{db: db}
}

// Doは指定された関数をデータベーストランザクション内で実行する
// 関数がエラーを返した場合、トランザクションはロールバックされる
// それ以外の場合、トランザクションはコミットされる
func (tm *txManager) Do(ctx context.Context, fn func(context.Context, domain.Executor) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// トランザクション用のexecutorラッパーを作成
	executor := NewTxExecutor(tx)

	// ビジネスロジックを実行
	if err := fn(ctx, executor); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction (original error: %w): %v", err, rbErr)
		}
		return err
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AsExecutorは非トランザクションのExecutorを返す
func (tm *txManager) AsExecutor() domain.Executor {
	return NewDBExecutor(tm.db)
}
