package radio

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"bytes"
	//"regexp"
	//"sync"
	"encoding/json"
	"time"
	//"google.golang.org/api/googleapi/transport"
	//"google.golang.org/api/youtube/v3"
)

type Jockey struct {
	newId   int
	queue   chan string
	mutex   chan int
	syncer  chan int
	Clients map[int]Client
	JukeBox JukeBox
}

func NewJockey(songList []string) *Jockey {
	mutex := make(chan int, 1)
	mutex <- 1

	return &Jockey{
		newId:   0,
		queue:   make(chan string),
		mutex:   mutex,
		syncer:  make(chan int),
		Clients: make(map[int]Client),
		JukeBox: JukeBox{songs: songList, position: 0, startTime: time.Now()}}
}

func (jockey *Jockey) Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Connection")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	newClient := Client{id: jockey.newId, endpoint: conn}

	closeChannel := make(chan int)
	go newClient.Listen(jockey.queue, closeChannel, jockey.syncer)
	go jockey.Disconnect(closeChannel)

	lock := <-jockey.mutex
	jockey.Clients[jockey.newId] = newClient
	jockey.newId ++
	jockey.mutex <- lock

	fmt.Println("ClientId - %d", jockey.newId)
}

func (jockey *Jockey) Disconnect(close chan int) {
	index := <-close
	lock := <-jockey.mutex
	delete(jockey.Clients, index)
	jockey.mutex <- lock
}

func (jockey *Jockey) PopulateQueue() {
	fmt.Println("Reading input channel ...")
	for {
		mediaId := <-jockey.queue
		fmt.Println("One of the clients entered - " + mediaId)

		jukebox := &(jockey.JukeBox)
		jukebox.songs = append(jukebox.songs, mediaId)

		fmt.Println(jockey.JukeBox.songs)

		msg, err := json.Marshal(Message{Kind: "Title", Value: mediaId})
		if err != nil {
			fmt.Println("Unable to form response JSON.")
			msg = []byte("Surprise Song!!")
		}

		for _, client := range jockey.Clients {
			client.endpoint.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (jockey *Jockey) Synchronize() {
	for {
		<-jockey.syncer

		jukeBox := &(jockey.JukeBox)
		msg := jukeBox.getNextSong()

		time.Sleep(5 * time.Second)

		fmt.Println("Sending next song to all clients %s", string(msg[:]))
		for _, client := range jockey.Clients {
			client.endpoint.WriteMessage(websocket.TextMessage, []byte(msg))
		}

		jukeBox.startTime = time.Now()
	}
}

func (jukeBox *JukeBox) getNextSong() []byte {
	if (jukeBox.position + 1) == len(jukeBox.songs) {
		msg, _ := json.Marshal(Message{Kind: "Finished"})
		return msg
	} else {
		jukeBox.position ++

		msg, _ := json.Marshal(Message{Kind: "Song", Value: jukeBox.songs[jukeBox.position]})
		return msg
	}
}


//func getTitle(id string) string {
//	client := &http.Client{
//		Transport: &transport.APIKey{Key: "AIzaSyCS-TiDxUSVGLuQfIyMHlhUdG9wiFu8d_A"},
//	}
//
//	service, err := youtube.New(client)
//	if err != nil {
//		fmt.Println("Error creating new YouTube client: %v", err)
//		return id
//	}
//
//	service.Sea
//}