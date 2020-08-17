package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

type sessionData struct {
	Error    string
	Saved    int
	Filename string
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

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err = tmpl.Execute(w, sessionData)
	if err != nil {
		panic(err)
	}
}

func showResult(w http.ResponseWriter, r *http.Request) {
	sessionData := sessionData{}

	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sessionValues := session.Flashes(); len(sessionValues) > 0 {
		sessionData.Filename = sessionValues[0].(string)
		sessionData.Saved = sessionValues[1].(int)
	}
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	err = tmpl.Execute(w, sessionData)
	if err != nil {
		panic(err)
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
