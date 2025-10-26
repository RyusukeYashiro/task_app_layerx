package clock

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// RealClock は domain.Clock の本番実装です
type RealClock struct{}

func New() domain.Clock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}
