package domain

import "context"

// TxManagerはデータベーストランザクションを管理
// ドメイン層から*sql.Txを隠蔽
type TxManager interface {
	Do(ctx context.Context, fn func(context.Context, Executor) error) error
	AsExecutor() Executor // 非トランザクションのExecutorを返す
}
