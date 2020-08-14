package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

type app struct {
	config config
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

	cfg := config{
		// at least 180 seconds to be able to handle big files uploads with slow 3G connection
		readWriteTimeout: 180,
		serverPort:       serverPort,
	}

	return &app{
		config: cfg,
	}
}

func (a *app) Run() {
	http.HandleFunc("/", showHomePage)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/result", showResult)
	http.HandleFunc("/images", download)
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	srv := &http.Server{
		Addr:         net.JoinHostPort("", a.config.serverPort),
		ReadTimeout:  a.config.readWriteTimeout * time.Second,
		WriteTimeout: a.config.readWriteTimeout * time.Second,
	}
	log.Println("Ready to start the server...")

	log.Println(srv.ListenAndServe())
}
