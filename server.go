package main

import "net/http"

func createRoutes() {
	http.HandleFunc("/", showHomePage)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/result", showResult)
	http.HandleFunc("/images", download)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
}
