package utils

import (
	"fmt"
	"path"
	"time"
)

func ValueOrDefaultString(data *map[string]string, key string, def string) string {
	v, exists := (*data)[key]

	if !exists {
		return def
	}

	return v
}

func ValueOrErrorString(data *map[string]string, key string, action string) (string, error) {
	v, exists := (*data)[key]

	if !exists {
		return "", fmt.Errorf("%s: Missing required field '%s'", action, key)
	}

	return v, nil
}

func AbortIfError(err interface{}) {
	if err != nil {
		Log.Fatal("%s", err)
	}
}

func AbortIf(cond bool, format string, v ...interface{}) {
	if cond {
		Log.Fatal(format, v...)
	}
}

type FileData struct {
	Name string
	Path string
}

func NewFileData(format, filePath, ext string) FileData {
	newFileName := fmt.Sprintf("%s.%s", FormatDate(format, time.Now()), ext)

	return FileData{
		Name: newFileName,
		Path: path.Join(filePath, newFileName),
	}
}