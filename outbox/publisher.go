package outbox

import "context"

type Publisher interface {
	Publish(context.Context, string, interface{}, ...map[string][]byte) error
}

type OutboxPublisher interface {
	PrePublish(Pre)
	StopProduce()
	Publisher
}

type Pre func(ctx context.Context, message *Message)
