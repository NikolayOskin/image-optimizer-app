package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

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

func showHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

	data.Errors = nil
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Requested filename is empty", 400)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	file, err := os.Open(imagesPath + fileName)
	if err != nil {
		http.Error(w, "File doesn't exists", 404)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Something goes wrong", 500)
		return
	}
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
