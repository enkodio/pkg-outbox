package client

import (
	"github.com/sirupsen/logrus"
	"gitlab.enkod.tech/pkg/transactionoutbox/migration/app"
	"gitlab.enkod.tech/pkg/transactionoutbox/pkg/logger"
)

func NewOutbox(
	pgClient Client,
	tx Transactor,
	publisher ReceivedPublisher,
	serviceName string,
	log *logrus.Logger,
) (RecordLogic, GivenPublisher) {
	if log != nil {
		logger.SetLogger(log)
	} else {
		logger.SetDefaultLogger("debug")
	}

	app.Run(pgClient.GetSqlDB(), serviceName, map[string]string{
		"c": "up",
	})

	var (
		storeRepository = newStoreRepository(pgClient)
	)

	var (
		recordLogic    = newRecordsLogic(storeRepository, tx, publisher)
		publisherLogic = newPublisherLogic(storeRepository, serviceName)
	)
	return recordLogic, publisherLogic
}
