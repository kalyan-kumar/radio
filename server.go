package main

import (
	"net/http"
	"html/template"
)

type pageParameters struct {

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, pageParameters{})
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":8042", nil)
}
