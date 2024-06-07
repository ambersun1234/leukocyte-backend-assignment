package orchestration

import (
	"leukocyte/src/types"
)

type Orchestration interface {
	Schedule(types.JobObject) error
}
