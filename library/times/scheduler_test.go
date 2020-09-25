package times

import (
	"startkit/library/systems"
	"testing"
	"time"
)

func TestParseAny(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		// TODO: Add test cases.
		{name: "A", value: "20160203"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseAny(tt.value)
		})
	}
}

func TestRoutine(t *testing.T) {
	type args struct {
		start    []int
		count    int
		interval time.Duration
		f        func() error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "a",
			args: args{
				start:    []int{0, 0, 0, 0},
				count:    10,
				interval: 2 * time.Second,
				f: func() error {
					var (
						err   error
						now   = time.Now().Local()
						upath = "files/"
						dpath = "files/"
						name  = ""
					)
					_, err = systems.MustOpen(name, upath)
					_, err = systems.MustOpen(name, dpath)
					_, err = systems.MustOpen(name, dpath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")+systems.GetSplit())
					_, err = systems.MustOpen(name, upath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")+systems.GetSplit())
					_, err = systems.MustOpen(name, dpath+systems.GetSplit()+time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")+systems.GetSplit()+"tn")
					return err
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Routine(tt.args.start, tt.args.count, tt.args.interval, tt.args.f)
		})
	}
}
