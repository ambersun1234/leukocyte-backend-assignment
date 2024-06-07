package orchestration

import (
	"leukocyte/src/types"
)

//go:generate mockery --name Orchestration
type Orchestration interface {
	Schedule(types.JobObject) error
}
