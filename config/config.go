package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Path string
	Hash string
	Actions actionsConfig
	Cryption cryptionConfig
}

type configFile struct {
	Path string `yaml:"path"`
	Hash string `yaml:"hash"`
	Actions map[string]map[string]string `yaml:"actions"`
	Encryption map[string]string `yaml:"encryption"`
}

func LoadConfig(cfgPath string) (*Config, error) {
	config, err := parseFile(cfgPath)
	if err != nil {
		return nil, err
	}

	cryption, err := loadCryption(config)
	if err != nil {
		return nil, err
	}

	actions, err := loadActions(config)
	if err != nil {
		return nil, err
	}

	return &Config{
		Actions: actions,
		Cryption: cryption,
	}, nil
}

func parseFile(cfgPath string) (*configFile, error) {
 
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("config/config: couldn't open config file")
	}
 
	parsed := &configFile{}
    
	err = yaml.Unmarshal([]byte(data), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}