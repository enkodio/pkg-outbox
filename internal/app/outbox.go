package app

import (
	"context"
	kafkaClient "github.com/enkodio/pkg-kafka/client"
	"github.com/enkodio/pkg-outbox/client"
	"github.com/enkodio/pkg-outbox/outbox"

	configEntity "github.com/enkodio/pkg-outbox/pkg/config/entity"
	"github.com/enkodio/pkg-outbox/pkg/logger"
	postgres "github.com/enkodio/pkg-postgres/client"
	"time"
)

func Run(configSettings configEntity.Settings, serviceName string) {
	logger.SetDefaultLogger("debug")
	const (
		testTopic = "test_topic"
	)
	var (
		pgClient, transactor = postgres.NewClient(configSettings.PostgresConfigs, serviceName, nil)
		k                    = kafkaClient.NewClient(
			configSettings.KafkaProducer,
			configSettings.KafkaConsumer,
			serviceName,
			nil,
			"local_")
		transactionOutbox, publisher = client.NewOutbox(pgClient, transactor, k, serviceName, nil)
	)
	testConsumer(testTopic, k)
	testProducer(testTopic, publisher)
	k.Start()

	transactionOutbox.Start()
	time.Sleep(time.Second * 10)
	transactionOutbox.Stop()
	select {}
}

func testConsumer(topic string, k kafkaClient.Client) {
	k.Subscribe(testHandler, 1, &kafkaClient.TopicSpecifications{
		NumPartitions:     1,
		ReplicationFactor: 1,
		Topic:             topic,
	})
}

func testProducer(topic string, k outbox.Publisher) {
	err := k.Publish(context.Background(), topic, "test", map[string][]byte{
		"test": []byte("test"),
	})
	if err != nil {
		logger.GetLogger().WithError(err).Error("produce err")
	}
}

func testHandler(ctx context.Context, msg []byte) error {
	logger.GetLogger().Info(string(msg))
	return nil
}
