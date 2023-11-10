package outbox

import "time"

type RecordLogic interface {
	Start()
	Stop()
	ChangeRecordSettings
}

type ChangeRecordSettings interface {
	SetSelectLimit(int)
	SetCountGoroutines(int)
	SetSleepTime(time.Duration)
}
