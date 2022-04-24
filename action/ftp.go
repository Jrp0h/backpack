package action

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Jrp0h/backpack/utils"
	"github.com/google/uuid"
	"github.com/jlaffaye/ftp"
)

type ftpAction struct {
	user     string
	password string
	host     string
	port     uint64
	dir      string
}

func (action *ftpAction) connect() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", action.host, action.port), ftp.DialWithTimeout(60*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(action.user, action.password)
	if err != nil {
		return nil, err
	}

	return c, nil

}

func (action *ftpAction) CanValidateConnection() bool {
	return true
}

func (action *ftpAction) TestConnection() error {
	c, err := action.connect()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}

func (action *ftpAction) Upload(fileData *utils.FileData) error {
	c, err := action.connect()
	if err != nil {
		return err
	}

	defer c.Quit()

	data, err := ioutil.ReadFile(fileData.Path)
	if err != nil {
		return err
	}

	err = c.ChangeDir(action.dir)
	if err != nil {
		return fmt.Errorf("couldn't change to %s. %s", action.dir, err.Error())
	}

	return c.Stor(fileData.Name, bytes.NewReader(data))
}

func (action *ftpAction) ListFiles() ([]string, error) {
	c, err := action.connect()
	if err != nil {
		return nil, err
	}
	defer c.Quit()

	entries, err := c.List(action.dir)
	if err != nil {
		return nil, fmt.Errorf("couldn't list %s. %s", action.dir, err.Error())
	}

	files := make([]string, 0)

	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			files = append(files, entry.Name)
		}
	}

	return files, nil
}

func (action *ftpAction) Fetch(file string) (string, error) {
	c, err := action.connect()
	if err != nil {
		return "", err
	}
	defer c.Quit()

	err = c.ChangeDir(action.dir)
	if err != nil {
		return "", fmt.Errorf("couldn't change to %s. %s", action.dir, err.Error())
	}

	r, err := c.Retr(file)
	if err != nil {
		return "", fmt.Errorf("couldn't retrive file %s. %s", file, err.Error())
	}
	defer r.Close()

	outPath := path.Join(os.TempDir(), uuid.NewString()+".zip")
	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, r); err != nil {
		return "", err
	}

	return outPath, nil
}

func loadFTPAction(data *map[string]string) (Action, error) {

	// Required
	user, err := utils.ValueOrErrorString(data, "user", "action/ftp")
	if err != nil {
		return nil, err
	}

	host, err := utils.ValueOrErrorString(data, "host", "action/ftp")
	if err != nil {
		return nil, err
	}

	// Optional
	password := utils.ValueOrDefaultString(data, "password", "")
	txtPort := utils.ValueOrDefaultString(data, "port", "21")
	port, err := strconv.ParseUint(txtPort, 10, 0)
	if err != nil || port > 65535 {
		return nil, fmt.Errorf("action/ftp: '%s' is not a valid port", txtPort)
	}
	dir := utils.ValueOrDefaultString(data, "dir", "/")

	return &ftpAction{
		user,
		password,
		host,
		port,
		dir,
	}, nil
}
