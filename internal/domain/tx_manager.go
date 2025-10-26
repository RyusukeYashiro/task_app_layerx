package domain

import "context"

// TxManager manages database transactions
// It hides *sql.Tx from domain layer
type TxManager interface {
	Do(ctx context.Context, fn func(context.Context, Executor) error) error
}
