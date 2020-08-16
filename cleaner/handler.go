package main

import (
	"io/ioutil"
	"os"
	"time"
)

func deleteOldImages(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := deleteIfOld(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteIfOld(file os.FileInfo) error {
	if !file.Mode().IsRegular() {
		return nil
	}
	if isOlderThanHour(file.ModTime()) {
		if err := os.Remove(imagesPath + file.Name()); err != nil {
			return err
		}
	}
	return nil
}

func isOlderThanHour(t time.Time) bool {
	return time.Now().Sub(t) > 1*time.Hour
}
