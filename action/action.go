package action

import (
	"fmt"
	"strings"

	"github.com/Jrp0h/backpack/utils"
)

type Action interface {
	CanValidateConnection() bool
	TestConnection() error
	Upload(*utils.FileData) error
	ListFiles() ([]string, error)
	Fetch(string) (string, error)
}

func LoadFromConfig(data *map[string]string) (Action, error) {
	t, exists := (*data)["type"]

	if !exists {
		return nil, fmt.Errorf("action/action: Missing required field 'type'")
	}

	switch strings.ToLower(t) {
	case "ftp":
		return loadFTPAction(data)
	case "move":
		return loadMoveAction(data)
	case "s3":
		return loadS3Action(data)
	default:
		return nil, fmt.Errorf("action/action: Unknown type '%s'", t)
	}
}
