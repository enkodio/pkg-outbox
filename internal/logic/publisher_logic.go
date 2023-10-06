package logic

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.enkod.tech/pkg/transactionoutbox/client"
	"gitlab.enkod.tech/pkg/transactionoutbox/internal/entity"
	"time"
)

type publisherLogic struct {
	serviceName     string
	storeRepository entity.Store
}

func NewPublisherLogic(storeRepository entity.Store, serviceName string) client.Publisher {
	return &publisherLogic{
		storeRepository: storeRepository,
		serviceName:     serviceName,
	}
}

func (p *publisherLogic) validateMessage(message entity.Message) error {
	if message.Topic == "" {
		return errors.New("invalid message, topic is empty")
	}
	if message.Body == nil {
		return errors.New("invalid message, body is empty")
	}
	return nil
}

func (p *publisherLogic) Publish(ctx context.Context, topic string, data interface{}, headers ...client.Header) error {
	message := entity.NewMessage(topic, data, headers)
	err := p.validateMessage(message)
	if err != nil {
		return err
	}
	record := entity.Record{
		ServiceName: p.serviceName,
		Uuid:        uuid.New(),
		Message:     message,
		State:       entity.PendingDelivery,
		CreatedOn:   time.Now().UTC().Unix(),
	}
	return p.storeRepository.AddRecord(ctx, record)
}
