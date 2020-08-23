package cleaner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func deleteOldFiles(dir string, olderThan time.Duration) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := deleteIfOld(file, dir, olderThan)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteIfOld(file os.FileInfo, dir string, olderThan time.Duration) error {
	if !file.Mode().IsRegular() {
		return nil
	}
	if time.Now().Sub(file.ModTime()) > olderThan {
		if err := os.Remove(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}
	return nil
}
