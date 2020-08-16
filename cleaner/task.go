package main

import (
	"log"
	"sync"
	"time"
)

type FileCleanerTask struct {
	closed chan struct{}
	ticker *time.Ticker
	wg     sync.WaitGroup
}

func (t *FileCleanerTask) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			err := deleteOldImages(imagesPath)
			if err != nil {
				log.Printf("error while delete old images: %v", err)
			}
		}
	}
}

func (t *FileCleanerTask) Stop() {
	close(t.closed)
	t.wg.Wait()
}
