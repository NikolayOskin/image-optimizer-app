package main

import (
	"github.com/gorilla/sessions"
)

const imagesPath string = "./images/"
const sessionName string = "tiny-images"

var store *sessions.CookieStore

func main() {
	app := NewApp()
	app.Run()
}
