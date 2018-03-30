package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/kalyan-kumar/radio/src"
)

func sockapiHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/web-socket.js")
	t.Execute(w, radio.PlayingSong{})
	fmt.Println("Sending sockapi")
}

/*
Syncing can instead be done by having song length in the server. It would give better experience
and solve the problem easily.
*/

func main() {
	jockey := radio.NewJockey([]radio.Song{})
	go jockey.PopulateQueue()
	go jockey.Synchronize()

	http.HandleFunc("/radio", jockey.JukeBox.LoadPage)
	http.HandleFunc("/ytapi", jockey.JukeBox.InitializePlayer)
	http.HandleFunc("/sockapi", sockapiHandler)
	http.HandleFunc("/websocket", jockey.Connect)

	http.ListenAndServe(":8042", nil)
}
