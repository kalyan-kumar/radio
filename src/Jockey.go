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

func (clients *Jockey) Connect(w http.ResponseWriter, r *http.Request) {
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

	newClient := Client{Endpoint: conn}

	closeChannel := make(chan int)
	go newClient.Listen(clients.queue, closeChannel, clients.syncer)
	go clients.Disconnect(closeChannel)

	lock := <-clients.mutex
	clients.Clients[clients.newId] = newClient
	clients.newId ++
	clients.mutex <- lock

	fmt.Println(clients.newId)
}

func (clients *Jockey) Disconnect(close chan int) {
	index := <-close
	lock := <-clients.mutex
	delete(clients.Clients, index)
	clients.mutex <- lock
}

func (clients *Jockey) PopulateQueue() {
	fmt.Println("Reading input channel ...")
	for {
		v := <-clients.queue
		fmt.Println("One of the clients entered - " + v)
		msg, err := json.Marshal(Message{Kind: "Title", Value: "v"})
		if err != nil {
			fmt.Println("Unable to form response JSON.")
			msg = []byte("Surprise Song!!")
		}

		for _, client := range clients.Clients {
			client.Endpoint.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (jockey *Jockey) Synchronize() {
	for {
		<-jockey.syncer

		jukeBox := jockey.JukeBox
		msg := jukeBox.getNextSong()

		fmt.Println("Sending next song to all clients")
		for _, client := range jockey.Clients {
			//dummy, _ := json.Marshal(Message{Kind: "Song", Value: "0z8wohG5mqI"})
			client.Endpoint.WriteMessage(websocket.TextMessage, []byte(msg))
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
