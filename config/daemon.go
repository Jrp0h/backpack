package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type TasksConfig struct {
	Config   string   `yaml:"config_path"`
	Interval string   `yaml:"interval"`
	Only     []string `yaml:"only"`
	Except   []string `yaml:"except"`
}

type tasksConfig map[string]TasksConfig

type DaemonConfig struct {
	Actions actionsConfig
	Tasks   tasksConfig
}

type daemonConfigFile struct {
	Actions map[string]map[string]string `yaml:"actions"`
	Tasks   map[string]TasksConfig       `yaml:"tasks"`
}

func LoadDaemonConfig(cfgPath string) (*DaemonConfig, error) {
	config, err := parseDaemonFile(cfgPath)
	if err != nil {
		return nil, err
	}

	actions, err := loadActions(&config.Actions)
	if err != nil {
		return nil, err
	}

	return &DaemonConfig{
		Actions: actions,
		Tasks:   config.Tasks,
	}, nil
}
func parseDaemonFile(cfgPath string) (*daemonConfigFile, error) {

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open config file")
	}

	parsed := &daemonConfigFile{}

	err = yaml.Unmarshal([]byte(data), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
