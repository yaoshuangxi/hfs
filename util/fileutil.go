package util

import (
	"os"
	"path/filepath"
)

func Mkdir(dir string) error {
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		if fi != nil && !fi.IsDir() {
			dir = filepath.Dir(dir)
			return Mkdir(dir)
		}
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}
