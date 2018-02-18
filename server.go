package main

import (
	//"fmt"
	"net/http"
	"html/template"
	//"github.com/gorilla/websocket"
)

type playingSong struct {
	Id       string
	Position int
}

type videoDetails struct {
	Id string
	Name string
}

type pageParameters struct {
	ToPlay playingSong
	Queue []videoDetails
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, pageParameters{Queue:[]videoDetails{{Name:"Bob Dylan"}, {Name:"Skylar Grey"}}})
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("js/youtube-api.js")
	t.Execute(w, pageParameters{ToPlay:playingSong{Id:"GDQob4AOCsQ", Position:30}})
}

func main() {
	//websocket.Upgrader{}
	http.HandleFunc("/radio", viewHandler)
	//http.Handle("/script", http.StripPrefix("/script", http.FileServer(http.Dir("js/youtube-api.js"))))
	http.HandleFunc("/script", scriptHandler)
	http.ListenAndServe(":8042", nil)
}
