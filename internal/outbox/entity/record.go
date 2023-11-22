package entity

import (
	"github.com/enkodio/pkg-outbox/outbox"
	"github.com/google/uuid"
)

// Record represents the record that is stored and retrieved from the database
type Record struct {
	ServiceName string
	Uuid        uuid.UUID
	Message     outbox.Message
	State       RecordState
	CreatedOn   int64
}

type Records []Record

func (r Records) GetUuids() []uuid.UUID {
	uuids := make([]uuid.UUID, len(r))
	for i := range r {
		uuids[i] = r[i].Uuid
	}
	return uuids
}

// RecordState is the State of the Record
type RecordState int

const (
	//PendingDelivery is the initial state of all records
	PendingDelivery RecordState = iota
	//Delivered indicates that the Records is already Delivered
	Delivered
	//DeliveredErr indicates that the message is not Delivered because an error occurred
	DeliveredErr
	//MaxAttemptsReached indicates that the message is not Delivered but the max attempts are reached so it shouldn't be delivered
	MaxAttemptsReached
)
