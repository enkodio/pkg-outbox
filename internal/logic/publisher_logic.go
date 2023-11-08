package logic

import (
	"context"
	"github.com/enkodio/pkg-outbox/internal/entity"
	"github.com/enkodio/pkg-outbox/outbox"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type PublisherLogic struct {
	serviceName     string
	storeRepository entity.Store
	prePublish      []outbox.Pre
}

func NewPublisherLogic(storeRepository entity.Store, serviceName string) *PublisherLogic {
	return &PublisherLogic{
		storeRepository: storeRepository,
		serviceName:     serviceName,
	}
}

func (p *PublisherLogic) validateMessage(message outbox.Message) error {
	if message.Topic == "" {
		return errors.New("invalid message, topic is empty")
	}
	if message.Body == nil {
		return errors.New("invalid message, body is empty")
	}
	return nil
}

func (p *PublisherLogic) Publish(ctx context.Context, topic string, data interface{}, headers ...map[string][]byte) error {
	message := outbox.NewMessage(topic, data, headers...)
	err := p.validateMessage(message)
	if err != nil {
		return err
	}
	for _, pre := range p.prePublish {
		pre(ctx, &message)
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

func (p *PublisherLogic) PrePublish(pre outbox.Pre) {
	p.prePublish = append(p.prePublish, pre)
}
