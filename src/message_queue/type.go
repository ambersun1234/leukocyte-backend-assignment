package message_queue

import (
	"leukocyte/src/types"
)

type Queue interface {
	Publish(string, string) error
	Consume(string, types.CallbackFunc) error
}
