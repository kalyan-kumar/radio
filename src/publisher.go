package radio

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"bytes"
	//"regexp"
	//"sync"
	"html/template"
	"encoding/json"
	"time"
)

const MaximumClients = 5

type Message struct {
	Kind  string
	Value string
}

type PlayingSong struct {
	Id       string
	Position int
}

type Client struct {
	Handle   string
	Id       int
	Endpoint *websocket.Conn
}

func (client Client) Listen(queue chan string, close chan int, syncer chan int) {
	fmt.Println("Running Listener ...")
	for {
		messageType, data, err := client.Endpoint.ReadMessage()
		if err != nil {
			fmt.Println("Error on reading socket. Closing the listener.")
			close <- client.Id
			break
		}
		fmt.Println(data)

		if messageType == websocket.TextMessage {
			inputString := string(data[:])
			fmt.Println("Received input - " + inputString)

			if "Finished" == inputString {
				select {
				case syncer <- 1:
					fmt.Println("Running Synchronizer")
				default:
					fmt.Println("Nothing to do")
				}
			} else {

				//m, _ := regexp.MatchString("^https://youtube\\.com/", inputString)
				//if m {
				queue <- inputString
			}
			//}
		}
	}
}

type SongList struct {
	songs     []string
	position  int
	startTime time.Time
}

type ClientList struct {
	newId    int
	queue    chan string
	mutex    chan int
	syncer   chan int
	Clients  map[int]Client
	songList SongList
}

func NewClientList() *ClientList {
	mutex := make(chan int, 1)
	mutex <- 1

	return &ClientList{
		newId:   0,
		queue:   make(chan string),
		mutex:   mutex,
		syncer:  make(chan int),
		Clients: make(map[int]Client)}
}

func (clients *ClientList) Connect(w http.ResponseWriter, r *http.Request) {
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

func (clients *ClientList) Disconnect(close chan int) {
	index := <-close
	lock := <-clients.mutex
	delete(clients.Clients, index)
	clients.mutex <- lock
}

func (clients *ClientList) PopulateQueue() {
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

func (clients *ClientList) InitializePlayer(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("js/youtube-api.js")

	songList := clients.songList
	id := songList.songs[songList.position]
	pos := int(time.Since(songList.startTime))

	t.Execute(w, PlayingSong{Id: id, Position: pos})
}

func (clients *ClientList) Synchronize() {
	for {
		<-clients.syncer

		songList := clients.songList

		if songList.position == len(songList.songs) {
			_, err := json.Marshal(Message{Kind: "Finished"})
			if err != nil {
				//msg = []byte("0z8wohG5mqI")
			}

		} else {
			clients.songList.position ++
			_, err := json.Marshal(Message{Kind: "Song", Value: songList.songs[songList.position]})
			if err != nil {
				fmt.Println("Unable to form response JSON for next song.")
				//msg = []byte("0z8wohG5mqI")
			}
		}

		fmt.Println("Sending next song to all clients")

		for _, client := range clients.Clients {
			dummy, _ := json.Marshal(Message{Kind: "Song", Value: "0z8wohG5mqI"})
			fmt.Println(dummy)
			client.Endpoint.WriteMessage(websocket.TextMessage, []byte(dummy))
		}

		songList.startTime = time.Now()
	}
}
