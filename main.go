package main

import (
	"fmt"
	"runtime"
)

const imagesPath string = "./images/"

type validationError struct {
	Error string
}

type appData struct {
	Errors           []validationError
	HandledImageName string
}

var data appData

const imagesPath string = "./images/"

func main() {
	data = appData{}
	http.HandleFunc("/", showHomePage)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/images", handleDownloadFile)
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	srv := &http.Server{
		Addr:         net.JoinHostPort("", "8083"),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
	}
	log.Println(srv.ListenAndServe())
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
