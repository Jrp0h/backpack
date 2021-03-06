package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Jrp0h/backpack/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Path           string
	Hash           string
	CWD            string
	FileNameFormat string

	Actions actionsConfig
	Crypto  cryptoConfig
}

const (
	Path int = 1 << iota
	Hash
	CWD

	Actions
	Crypto
)

func (c *Config) Require(fields int) {
	if fields&Path != 0 {
		checkStringField("data_path", c.Path)
	}

	if fields&Hash != 0 {
		checkStringField("hash_path", c.Hash)
	}

	if fields&CWD != 0 {
		checkStringField("cwd", c.CWD)
	}

	if fields&Actions != 0 {
		if len(c.Actions) < 1 {
			utils.Log.Fatal("config/config: atleast one action is required")
		}
	}

	if fields&Crypto != 0 {
		if !c.Crypto.Enabled {
			utils.Log.Fatal("config/config: encryption is required")
		}
	}
}

func (c *Config) Validate(fields int) {
	if fields&Path != 0 {
		if !utils.PathExists(c.Path) {
			utils.Log.Fatal("config/config: data_path isn't a valid path")
		}
	}

	if fields&Hash != 0 {
		if !utils.PathExists(c.Hash) || !utils.PathIsFile(c.Hash) {
			utils.Log.Fatal("config/config: hash_path isn't a valid path or isn't a file")
		}
	}

	if fields&CWD != 0 {
		if !utils.PathExists(c.Hash) {
			utils.Log.Fatal("config/config: cwd isn't a valid path")
		}
	}
}

func (c *Config) Cd() {
	if c.CWD != "" {
		if err := os.Chdir(c.CWD); err != nil {
			utils.Log.Fatal("%s", err.Error())
		}
	}
}

func checkStringField(name, value string) {
	if value == "" {
		utils.Log.Fatal("config/config: missing required option '%s'", name)
	}
}

type configFile struct {
	Path           string `yaml:"data_path"`
	Hash           string `yaml:"hash_path"`
	CWD            string `yaml:"cwd"`
	FileNameFormat string `yaml:"file_name_format"`

	Actions    map[string]map[string]string `yaml:"actions"`
	Encryption map[string]string            `yaml:"encryption"`
}

func LoadConfig(cfgPath string) (*Config, error) {
	config, err := parseFile(cfgPath)
	if err != nil {
		return nil, err
	}

	crypto, err := loadCrypto(config)
	if err != nil {
		return nil, err
	}

	actions, err := loadActions(config)
	if err != nil {
		return nil, err
	}

	if config.FileNameFormat == "" {
		config.FileNameFormat = "%Y-%m-%d_%H%M%S"
	}

	return &Config{
		Path:           config.Path,
		Hash:           config.Hash,
		CWD:            config.CWD,
		FileNameFormat: config.FileNameFormat,

		Actions: actions,
		Crypto:  crypto,
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
