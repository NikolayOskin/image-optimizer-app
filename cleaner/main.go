package cleaner

import (
	"time"
)

// NewTask is the FileCleanerTask constructor
func NewTask(dir string, interval time.Duration, olderThan time.Duration) *FileCleanerTask {
	return &FileCleanerTask{
		olderThan: olderThan,
		path:      dir,
		closed:    make(chan struct{}),
		ticker:    time.NewTicker(interval),
	}
}
