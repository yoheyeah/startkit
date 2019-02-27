package times

import (
	"time"
)

func Zero() time.Time {
	return time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
}
