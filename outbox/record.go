package outbox

import "time"

type RecordLogic interface {
	StartProcessRecords()
	StopProcessRecords()
	ChangeRecordSettings
}

type ChangeRecordSettings interface {
	SetSelectLimit(int)
	SetCountGoroutines(int)
	SetSleepTime(time.Duration)
}
