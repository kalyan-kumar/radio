package main

import (
	//"fmt"
	"net/http"
	"html/template"
	"github.com/kalyan-kumar/radio/src"
)

type playingSong struct {
	Id       string
	Position int
}

type videoDetails struct {
	Id   string
	Name string
}

type pageParameters struct {
	ToPlay playingSong
	Queue  []videoDetails
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, pageParameters{Queue: []videoDetails{{Name: "Bob Dylan"}, {Name: "Skylar Grey"}}})
}

func ytapiHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/youtube-api.js")
	t.Execute(w, pageParameters{ToPlay: playingSong{Id: "hCQhRDvayos", Position: 30}})
}

func sockapiHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/web-socket.js")
	t.Execute(w, pageParameters{})
}

func main() {
	//http.Handle("/script", http.StripPrefix("/script", http.FileServer(http.Dir("js/youtube-api.js"))))

	songList := []string{}
	clients := &radio.ClientList{NumOfClients: 0, Queue: make(chan string)}
	go clients.PopulateQueue(&songList)

	http.HandleFunc("/radio", viewHandler)
	http.HandleFunc("/ytapi", ytapiHandler)
	http.HandleFunc("/sockapi", sockapiHandler)
	http.HandleFunc("/websocket", clients.Connect)

	http.ListenAndServe(":8042", nil)
}
