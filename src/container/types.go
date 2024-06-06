package container

import (
	"leukocyte/src/types"
)

type Container interface {
	Schedule(types.JobObject) error
}
