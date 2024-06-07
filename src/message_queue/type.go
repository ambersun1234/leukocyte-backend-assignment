package message_queue

import (
	"leukocyte/src/types"
)

//go:generate mockery --name Queue
type Queue interface {
	Publish(string, string) error
	Consume(string, types.CallbackFunc) error
	Connect() error
	Close() error
}
