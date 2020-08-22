package cleaner

import (
	"log"
	"sync"
	"time"
)

// FileCleanerTask is the time interval task that checks and delete files
// which were created more than provided duration time ago.
type FileCleanerTask struct {
	olderThan time.Duration
	path      string
	closed    chan struct{}
	ticker    *time.Ticker
	Wg        sync.WaitGroup
}

func (t *FileCleanerTask) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			err := deleteOldFiles(t.path, t.olderThan)
			if err != nil {
				log.Printf("error while delete old images: %v", err)
			}
		}
	}
}

func (t *FileCleanerTask) Stop() {
	close(t.closed)
	t.Wg.Wait()
}
