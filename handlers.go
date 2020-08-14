package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

type request struct {
	Errors []validationError
}

const requestCtx = "requestCtx"

func showHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	ctx := r.Context().Value(requestCtx).(request)
	err := tmpl.Execute(w, ctx)
	if err != nil {
		panic(err)
	}

	//data.Errors = nil
}

func showResult(w http.ResponseWriter, r *http.Request) {

}

func download(w http.ResponseWriter, r *http.Request) {
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
