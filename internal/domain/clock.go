package domain

import "time"

// Clock is an interface for time operations (mockable for testing)
type Clock interface {
	Now() time.Time
}
