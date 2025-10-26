package domain

import "time"

// Clockは時刻操作のインターフェース（テスト用にモック可能）
type Clock interface {
	Now() time.Time
}
