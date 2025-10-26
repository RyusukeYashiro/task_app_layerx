package domain

import "context"

// Executor is an abstraction for query execution (works for both *sql.DB and *sql.Tx)
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) Row
}

// Result is an abstraction of sql.Result
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Rows is an abstraction of *sql.Rows
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row is an abstraction of *sql.Row
type Row interface {
	Scan(dest ...any) error
}
