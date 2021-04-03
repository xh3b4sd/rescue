package rescue

import (
	"github.com/xh3b4sd/rescue/pkg/task"
)

type Interface interface {
	Create(t *task.Task) error
	Delete(t *task.Task) error
	Expire() error
	Search() (*task.Task, error)
}
