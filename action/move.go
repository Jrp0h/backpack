package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Jrp0h/backpack/utils"
	"github.com/google/uuid"
)

type moveAction struct {
	dir string
}

func (action *moveAction) CanValidateConnection() bool {
	return true
}

// TODO: Maybe check if directory exists and if user has permissions
func (action *moveAction) TestConnection() error {
	if utils.PathIsDir(action.dir) {
		return nil
	}

	return fmt.Errorf("action/move: %s isn't a directory", action.dir)
}

func (action *moveAction) Upload(fileData *utils.FileData) error {
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

func (action *moveAction) ListFiles() ([]string, error) {
	files := make([]string, 0)
	entries, err := os.ReadDir(action.dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

func (action *moveAction) Fetch(file string) (string, error) {

	outputPath := path.Join(os.TempDir(), uuid.NewString()+".zip")
	if utils.PathExists(outputPath) {
		return "", fmt.Errorf("action/move: output path %s already exists", outputPath)
	}

	inputPath := path.Join(action.dir, file)

	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("action/move: couldn't read file %s\n%s", inputPath, err.Error())
	}

	err = ioutil.WriteFile(outputPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("action/move: couldn't write file %s\n%s", outputPath, err.Error())
	}

	return outputPath, nil
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
