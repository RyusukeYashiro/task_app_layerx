package mysql

import (
	"context"
	"database/sql"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// DBExecutor wraps *sql.DB to implement domain.Executor
type DBExecutor struct {
	db *sql.DB
}

func NewDBExecutor(db *sql.DB) domain.Executor {
	return &DBExecutor{db: db}
}

func (e *DBExecutor) ExecContext(ctx context.Context, query string, args ...any) (domain.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

func (e *DBExecutor) QueryContext(ctx context.Context, query string, args ...any) (domain.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

func (e *DBExecutor) QueryRowContext(ctx context.Context, query string, args ...any) domain.Row {
	return e.db.QueryRowContext(ctx, query, args...)
}

// TxExecutor wraps *sql.Tx to implement domain.Executor
type TxExecutor struct {
	tx *sql.Tx
}

func NewTxExecutor(tx *sql.Tx) domain.Executor {
	return &TxExecutor{tx: tx}
}

func (e *TxExecutor) ExecContext(ctx context.Context, query string, args ...any) (domain.Result, error) {
	return e.tx.ExecContext(ctx, query, args...)
}

func (e *TxExecutor) QueryContext(ctx context.Context, query string, args ...any) (domain.Rows, error) {
	return e.tx.QueryContext(ctx, query, args...)
}

func (e *TxExecutor) QueryRowContext(ctx context.Context, query string, args ...any) domain.Row {
	return e.tx.QueryRowContext(ctx, query, args...)
}
