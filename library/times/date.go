package times

import (
	"time"
)

func DateString() string {
	var (
		start = []int{0, 0, 0, 0}
		now   = time.Now().Local()
	)
	return time.Date(now.Year(), now.Month(), now.Day(), start[0], start[1], start[2], start[3], time.Local).Format("2006-01-02")
}
