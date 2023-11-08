package client

import (
	"github.com/enkodio/pkg-outbox/internal/logic"
	"github.com/enkodio/pkg-outbox/internal/repository"
	"github.com/enkodio/pkg-outbox/migration/app"
	"github.com/enkodio/pkg-outbox/outbox"
	"github.com/enkodio/pkg-outbox/pkg/logger"
	"github.com/sirupsen/logrus"
)

func NewOutbox(
	pgClient outbox.RepositoryClient,
	tx outbox.Transactor,
	publisher outbox.Publisher,
	serviceName string,
	log *logrus.Logger,
) (outbox.RecordLogic, outbox.OutboxPublisher) {
	if log != nil {
		logger.SetLogger(log)
	} else {
		logger.SetDefaultLogger("debug")
	}

	app.Run(pgClient.GetSqlDB(), serviceName, map[string]string{
		"c":   "up",
		"dir": "/migration/migrations",
	})

	var (
		storeRepository = repository.NewStoreRepository(pgClient)
	)

	var (
		recordLogic    = logic.NewRecordsLogic(storeRepository, tx, publisher)
		publisherLogic = logic.NewPublisherLogic(storeRepository, serviceName)
	)
	return recordLogic, publisherLogic
}
