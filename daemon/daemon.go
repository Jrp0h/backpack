package daemon

import (
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/task"
)

type Daemon struct {
	tasks []*task.Task
}

func New(cfgPath string) (*Daemon, error) {
	cfg, err := config.LoadDaemonConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	tasks, err := task.LoadDaemonTasks(cfg)
	if err != nil {
		return nil, err
	}

	return &Daemon{
		tasks: tasks,
	}, nil
}

func (d Daemon) Run() {

}
