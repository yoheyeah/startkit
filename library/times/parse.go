package times

import (
	"time"
)

func TimeParse(str, layout string) (t time.Time, err error) {
	if layout == "" {
		layout = "2006-01-02"
	}
	t, err = time.Parse(layout, str)
	return
}
