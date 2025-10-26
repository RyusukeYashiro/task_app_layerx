package clock_test

import (
	"testing"
	"time"

	"github.com/ryusuke/task_app_layerx/internal/infrastructure/clock"
)

func TestRealClock_Now(t *testing.T) {
	c := clock.New()

	before := time.Now()
	now := c.Now()
	after := time.Now()

	// 現在時刻が取得できること
	if now.Before(before) || now.After(after) {
		t.Errorf("時刻が範囲外です: before=%v, now=%v, after=%v", before, now, after)
	}
}

func TestRealClock_Multiple(t *testing.T) {
	c := clock.New()

	first := c.Now()
	time.Sleep(10 * time.Millisecond)
	second := c.Now()

	// 時間が進んでいること
	if !second.After(first) {
		t.Errorf("時間が進んでいません: first=%v, second=%v", first, second)
	}
}
