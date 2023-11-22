package entity

import (
	"time"
)

type RecordSettings struct {
	SelectLimit     int
	CountGoroutines int
	SleepTime       time.Duration
}

func (r *RecordSettings) SetSelectLimit(limit int) {
	r.SelectLimit = limit
}

func (r *RecordSettings) SetCountGoroutines(countGoroutines int) {
	r.CountGoroutines = countGoroutines
}

func (r *RecordSettings) SetSleepTime(duration time.Duration) {
	r.SleepTime = duration
}
