package entity

import (
	"context"
)

// Store is the interface that should be implemented by SQL-like database drivers to support the outbox functionality
type Store interface {
	//AddRecordTx stores the message within the provided database transaction
	AddRecord(context.Context, Record) error
	GetPendingRecords(context.Context, Filter) ([]Record, error)
	UpdateRecordsStatus(context.Context, Records, RecordState) error
	DeleteRecords(context.Context, Records) error
}
