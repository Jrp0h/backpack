package utils

import "os"

func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func PathIsFile(path string) bool  {
	if f, err := os.Stat(path); err == nil {
		return !f.IsDir()
	}

	return false
}

func PathIsDir(path string) bool  {
	if f, err := os.Stat(path); err == nil {
		return f.IsDir()
	}

	return false
}