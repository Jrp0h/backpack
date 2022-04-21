package action

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Jrp0h/backuper/utils"
	"github.com/jlaffaye/ftp"
)

type ftpAction struct {
	user string
	password string
	host string
	port uint64
	dir string
}

func (action *ftpAction) TestConnection() error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", action.host, action.port), ftp.DialWithTimeout(60*time.Second))
	if err != nil {
		return err
	}
	defer c.Quit()

	err = c.Login(action.user, action.password)
	if err != nil {
		return err
	}

	return nil
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