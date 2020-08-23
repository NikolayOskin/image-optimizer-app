package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type sessionData struct {
	Error      string
	Saved      int64
	Filename   string
	BeforeSize string
	AfterSize  string
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	sessionData := sessionData{}

	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errors := session.Flashes(); len(errors) > 0 {
		sessionData.Error = errors[0].(string)
	}
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "home.html")))
	err = tmpl.Execute(w, sessionData)
	if err != nil {
		panic(err)
	}
}

func showResult(w http.ResponseWriter, r *http.Request) {
	var (
		ok          bool
		sessionData sessionData
	)

	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionData.Filename, ok = session.Values["filename"].(string)
	sessionData.BeforeSize, ok = session.Values["beforeSize"].(string)
	sessionData.AfterSize, ok = session.Values["afterSize"].(string)
	sessionData.Saved, ok = session.Values["saved"].(int64)

	if !ok {
		http.Redirect(w, r, "/", 301)
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "result.html")))
	err = tmpl.Execute(w, sessionData)
	if err != nil {
		http.Error(w, "Something goes wrong", 500)
		return
	}
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
