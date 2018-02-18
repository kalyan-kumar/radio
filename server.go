package main

import (
	//"fmt"
	"net/http"
	"html/template"
	//"github.com/gorilla/websocket"
)

type pageParameters struct {

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, pageParameters{})
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("js/youtube-api.js")
	t.Execute(w, pageParameters{})
}

func main() {
	//websocket.Upgrader{}
	http.HandleFunc("/radio", viewHandler)
	//http.Handle("/script", http.StripPrefix("/script", http.FileServer(http.Dir("js/youtube-api.js"))))
	http.HandleFunc("/script", scriptHandler)
	http.ListenAndServe(":8042", nil)
}
