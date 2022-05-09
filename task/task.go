package task

import (
	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
)

type Task struct {
	cfg *config.Config
}

func NewTask(path string) (*Task, error) {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return &Task{
		cfg: cfg,
	}, nil
}

func LoadDaemonTasks(cfg *config.DaemonConfig) ([]*Task, error) {
	tasks := make([]*Task, 0)

	for tk, tv := range cfg.Tasks {
		t, err := NewTask(tv.Config)
		if err != nil {
			return nil, err
		}
		for ak, av := range cfg.Actions {
			_, exists := t.cfg.Actions[ak]
			if exists {
				utils.Log.Warning("Action %s already exist on task %s. Will override.", ak, tk)
			}

			t.cfg.Actions[ak] = av
		}

		if len(tv.Only) > 0 {
			a, err := t.cfg.Actions.Only(tv.Only)
			if err != nil {
				return nil, err
			}

			t.cfg.Actions = a
		}

		if len(tv.Except) > 0 {
			a, err := t.cfg.Actions.Except(tv.Except)
			if err != nil {
				return nil, err
			}

			t.cfg.Actions = a
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}
