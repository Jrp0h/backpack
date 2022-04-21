package action

import "github.com/Jrp0h/backuper/utils"

type moveAction struct {
	dir string
}

// TODO: Maybe check if directory exists and if user has permissions
func (action *moveAction) TestConnection() error {
	return nil
}

func loadMoveAction(data *map[string]string) (Action, error) {

	// Required
	dir, err := utils.ValueOrErrorString(data, "dir", "action/move")
	if err != nil {
		return nil, err
	}

	return &moveAction{
		dir,
	}, nil
}