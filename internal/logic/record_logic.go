package logic

import (
	"context"
	"github.com/enkodio/pkg-outbox/internal/entity"
	"github.com/enkodio/pkg-outbox/outbox"
	"github.com/enkodio/pkg-outbox/pkg/logger"
	pgClient "github.com/enkodio/pkg-postgres/client"
	"github.com/pkg/errors"
	"time"
)

const (
	defaultLimit                   = 100
	defaultCountGoroutines         = 1
	defaultProcessRecordsSleepTime = time.Second * 5
)

type RecordsLogic struct {
	storeRepository entity.Store
	transactor      outbox.Transactor
	broker          outbox.Publisher
	entity.RecordSettings
	syncGroup *entity.SyncGroup
}

func NewRecordsLogic(
	storeRepository entity.Store,
	transactor pgClient.Transactor,
	broker outbox.Publisher,
) *RecordsLogic {
	r := &RecordsLogic{
		storeRepository: storeRepository,
		transactor:      transactor,
		broker:          broker,
		syncGroup:       entity.NewSyncGroup(),
		RecordSettings: entity.RecordSettings{
			SelectLimit:     defaultLimit,
			CountGoroutines: defaultCountGoroutines,
			SleepTime:       defaultProcessRecordsSleepTime,
		},
	}
	return r
}

func (r *RecordsLogic) Start() {
	r.syncGroup.Add(r.CountGoroutines)
	for i := 0; i < r.CountGoroutines; i++ {
		go r.processRecords()
	}
}

func (r *RecordsLogic) Stop() {
	r.syncGroup.Close()
}

func (r *RecordsLogic) processRecords() {
	defer r.syncGroup.Done()
	ctx := context.Background()
	log := logger.GetLogger()
	for {
		time.Sleep(r.SleepTime)
		select {
		case <-r.syncGroup.IsDone():
			return
		default:
			err := r.processRecordsWork(ctx)
			if err != nil {
				log.WithError(err).Error("cant process records")
			}
		}
	}
}

func (r *RecordsLogic) processRecordsWork(ctx context.Context) error {
	err := r.transactor.Begin(&ctx)
	if err != nil {
		return errors.Wrap(err, "cant begin tx for process records")
	}
	defer r.transactor.Rollback(&ctx)
	records, err := r.storeRepository.GetPendingRecords(ctx, entity.Filter{
		Limit: r.SelectLimit,
	})
	if err != nil {
		return errors.Wrap(err, "cant get pending records")
	}
	if len(records) == 0 {
		return nil
	}

	log := logger.GetLogger()
	successfulRecords, errorRecords := r.publishRecords(ctx, records)
	//TODO maybe we shouldn’t update the records status to err and stay pending
	if len(errorRecords) > 0 {
		err = r.storeRepository.UpdateRecordsStatus(ctx, errorRecords, entity.DeliveredErr)
		if err != nil {
			log.WithError(err).Error("cant update status for records with error")
		}
	}
	if len(successfulRecords) > 0 {
		err = r.storeRepository.DeleteRecords(ctx, successfulRecords)
		if err != nil {
			log.WithError(err).Error("cant delete successful delivered records")
		}
	}

	return r.transactor.Commit(&ctx)
}

func (r *RecordsLogic) publishRecords(ctx context.Context, records []entity.Record) (entity.Records, entity.Records) {
	successfulRecords := make([]entity.Record, 0, len(records))
	errorRecords := make([]entity.Record, 0, len(records))
	for _, record := range records {
		err := r.broker.Publish(ctx, record.Message.Topic, record.Message.Body, record.Message.Headers.ToMap())
		if err != nil {
			errorRecords = append(errorRecords, record)
		}
		successfulRecords = append(successfulRecords, record)
	}
	return successfulRecords, errorRecords
}
