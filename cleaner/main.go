package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const imagesPath string = "../images/"

func main() {
	task := &FileCleanerTask{
		closed: make(chan struct{}),
		ticker: time.NewTicker(1 * time.Minute), // run task once per minute
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	task.wg.Add(1)
	go func() {
		defer task.wg.Done()
		task.Run()
	}()

	select {
	case sig := <-shutdown:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
	}
}
