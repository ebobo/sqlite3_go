package main

import (
	"os"
)

func makeDirIfNotExists(dirpath string) error {
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.Mkdir(dirpath, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
