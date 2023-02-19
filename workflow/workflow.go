package workflow

import "log"

type Workflower interface {
	Run(*log.Logger) error
}
