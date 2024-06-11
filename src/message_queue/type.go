package message_queue

import (
	"leukocyte/src/types"
)

//go:generate mockery --name Queue
type Queue interface {
	Publish(types.RoutingKey, string) error
	Consume(types.RoutingKey, types.CallbackFunc) error
	Connect() error
	Close() error
}
