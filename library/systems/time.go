package systems

import (
	"time"
)

func NowInUNIX() int64 {
	return time.Now().Unix()
}
