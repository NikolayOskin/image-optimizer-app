package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type validationError struct {
	Error string
}

type appData struct {
	Errors           []validationError
	HandledImageName string
}

var data appData

func main() {
	data = appData{}
	http.HandleFunc("/", showHomePage)
	http.HandleFunc("/upload", handleUploadedForm)
	http.HandleFunc("/images", handleDownloadFile)
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, data)
	data.Errors = nil
}

func handleUploadedForm(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1000000000)
	file, header, err := r.FormFile("file")
	switch err {
	case nil:
	case http.ErrMissingFile:
		data.Errors = append(data.Errors, validationError{Error: "You didn't choose file to upload"})
		log.Println("no file")
		http.Redirect(w, r, "/", 301)
		return
	default:
		log.Println(err)
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		panic(err.Error())
	}
	// set position back to start.
	if _, err := file.Seek(0, 0); err != nil {
		panic(err.Error())
	}
	if isFiletypeValid(fileHeader) == false {
		http.Redirect(w, r, "/", 301)
		return
	}
	uploadImage(&file, header.Filename)
	http.Redirect(w, r, "/", 301)
	return
}

func uploadImage(file *multipart.File, filename string) {
	uploadedFile, err := os.OpenFile("./images/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer uploadedFile.Close()
	_, err = io.Copy(uploadedFile, *file)
	if err != nil {
		panic(err.Error())
	}
}

func isFiletypeValid(fileHeader []byte) bool {
	contentType := http.DetectContentType(fileHeader)
	if contentType != "image/jpeg" && contentType != "image/jpg" {
		data.Errors = append(data.Errors, validationError{Error: "Filetype is not correct. Allowed types: jpg, jpeg, png, gif"})
		return false
	}
	return true
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	requestedFileName := r.URL.Query().Get("file")
	if requestedFileName == "" {
		http.Error(w, "Requested filename is empty", 400)
	}
	file, err := os.Open("./images/" + requestedFileName)
	if err != nil {
		http.Error(w, "File doesn't exists", 404)
	}
	defer file.Close()
	fmt.Println(file.Stat())
	w.Header().Set("Content-Disposition", "attachment; filename="+requestedFileName)
	io.Copy(w, file)
	return
}
