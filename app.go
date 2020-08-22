package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikolayOskin/image-optimizer-app/cleaner"
	"github.com/gorilla/sessions"
)

type app struct {
	config config
	task   *cleaner.FileCleanerTask
}

type config struct {
	serverPort       string
	readWriteTimeout time.Duration
}

func NewApp() *app {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Fatal("server port is empty")
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("session key is not set")
	}

	store = sessions.NewCookieStore([]byte(sessionKey))
	store.Options.MaxAge = 3600

	cfg := config{
		// at least 180 seconds to be able to handle big files uploads with slow 3G connection
		readWriteTimeout: 180,
		serverPort:       serverPort,
	}

	return &app{
		config: cfg,
		task:   cleaner.NewTask(imagesPath, 1*time.Minute, 5*time.Minute),
	}
}

func (a *app) Run() {
	createRoutes()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         net.JoinHostPort("", a.config.serverPort),
		ReadTimeout:  a.config.readWriteTimeout * time.Second,
		WriteTimeout: a.config.readWriteTimeout * time.Second,
	}

	go func() {
		log.Printf("API listening on %s", net.JoinHostPort("", a.config.serverPort))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %v", err)
		}
	}()

	log.Println("Starting file cleaner task")
	a.task.Wg.Add(1)
	go func() {
		a.task.Run()
		a.task.Wg.Done()
	}()

	select {
	case sig := <-shutdown:
		log.Printf("Got %s signal. Stopping cleaner task...\n", sig)
		a.task.Stop()
		log.Println("Shutting down the server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		_ = srv.Close()
		log.Println("Could not stop the server gracefully")
		return
	}
	log.Println("Server stopped")
}
