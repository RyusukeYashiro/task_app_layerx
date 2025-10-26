package clock

import (
	"time"

	"github.com/ryusuke/task_app_layerx/internal/domain"
)

// RealClockはdomain.Clockの実装
type RealClock struct{}

func New() domain.Clock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}
