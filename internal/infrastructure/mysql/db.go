package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Configは接続設定を保持する
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewDBで新しいMySQLデータベース接続を作成
func NewDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// コネクションプール設定
	db.SetMaxOpenConns(25)                 // 最大オープン接続数
	db.SetMaxIdleConns(10)                 // 最大アイドル接続数
	db.SetConnMaxLifetime(5 * time.Minute) // 最大接続ライフタイム

	// 接続を検証
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
