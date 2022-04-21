package action

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Jrp0h/backuper/utils"
)

type moveAction struct {
	dir string
}

// TODO: Maybe check if directory exists and if user has permissions
func (action *moveAction) TestConnection() error {
	if utils.PathIsDir(action.dir) {
		return nil
	}

	return fmt.Errorf("action/move: %s isn't a directory", action.dir)
}

func (action *moveAction) Run(fileData *utils.FileData) error {
	err := action.TestConnection()
	if err != nil {
		return err
	}

	outputPath := path.Join(action.dir, fileData.Name)
	if utils.PathExists(outputPath) {
		return fmt.Errorf("action/move: output path %s already exists", outputPath)
	}

	data, err := ioutil.ReadFile(fileData.Path)
    if err != nil {
		return fmt.Errorf("action/move: couldn't read file %s\n%s", fileData.Path, err.Error())
    }

    err = ioutil.WriteFile(outputPath, data, 0644)
    if err != nil {
		return fmt.Errorf("action/move: couldn't write file %s\n%s", outputPath, err.Error())
    }

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