package action

import (
	"fmt"
	"strings"
)

type Action interface {
	TestConnection() error
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
	default:
		return nil, fmt.Errorf("action/action: Unknown type '%s'", t)
	}
}
