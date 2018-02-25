package radio

import (
	"time"
	"net/http"
	"html/template"
	"fmt"
)

type JukeBox struct {
	songs     []string
	position  int
	startTime time.Time
}

func (jukeBox *JukeBox) InitializePlayer(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/youtube-api.js")

	id := jukeBox.songs[jukeBox.position]
	pos := int(time.Since(jukeBox.startTime) / time.Second)
	fmt.Println("Sending pos - %d", pos)

	t.Execute(w, PlayingSong{Id: id, Position: pos})
}
