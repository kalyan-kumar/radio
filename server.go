package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/kalyan-kumar/radio/src"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, radio.PlayingSong{})
	fmt.Println("View sent")
}

func ytapiHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/youtube-api.js")
	t.Execute(w, radio.PlayingSong{Id: "hCQhRDvayos", Position: 150})
	fmt.Println("Sending ytapi")
}

func sockapiHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/web-socket.js")
	t.Execute(w, radio.PlayingSong{})
	fmt.Println("Sending sockapi")
}

func main() {
	//http.Handle("/script", http.StripPrefix("/script", http.FileServer(http.Dir("js/youtube-api.js"))))

	clients := radio.NewClientList()
	go clients.PopulateQueue()
	go clients.Synchronize()

	http.HandleFunc("/radio", viewHandler)
	http.HandleFunc("/ytapi", ytapiHandler)
	http.HandleFunc("/sockapi", sockapiHandler)
	http.HandleFunc("/websocket", clients.Connect)

	http.ListenAndServe(":8042", nil)
}
