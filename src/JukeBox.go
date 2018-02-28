package radio

import (
	"time"
	"net/http"
	"html/template"
	"fmt"
)

type Song struct {
	Id string
	Title string
	Image string
}

/*
Extract out the channel data into a synchronizer data structure.
*/

type JukeBox struct {
	songs     []Song
	position  int
	startTime time.Time
	ended     bool
	mutex     chan bool
}

func (jukeBox *JukeBox) LoadPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, PlayList{SongList: jukeBox.songs})
}

func (jukeBox *JukeBox) InitializePlayer(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/youtube-api.js")

	if jukeBox.position == -1 {
		t.Execute(w, PlayingSong{})
		return
	}

	id := jukeBox.songs[jukeBox.position].Id
	pos := int(time.Since(jukeBox.startTime) / time.Second)
	fmt.Println("Sending pos - %d", pos)

	t.Execute(w, PlayingSong{Id: id, Position: pos})
}
