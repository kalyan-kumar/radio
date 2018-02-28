package radio

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"time"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

/*
Extract out the channel data into a synchronizer data structure.
*/

type Jockey struct {
	newId   int
	queue   chan string
	mutex   chan int
	syncer  chan int
	Clients map[int]Client
	JukeBox JukeBox
}

func NewJockey(songList []Song) *Jockey {
	mutex := make(chan int, 1)
	mutex <- 1

	jukeBoxMutex := make(chan bool, 1)
	jukeBoxMutex <- false

	length := len(songList)

	return &Jockey{
		newId:   0,
		queue:   make(chan string),
		mutex:   mutex,
		syncer:  make(chan int),
		Clients: make(map[int]Client),
		JukeBox: JukeBox{songs: songList, position: length - 1, startTime: time.Now(), ended: length == 0, mutex: jukeBoxMutex}}
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
		query := <-jockey.queue
		fmt.Println("API query string - " + query)

		id, title, image := getIdAndTitle(query)
		if "" == id {
			continue
		}

		song := Song{Id: id, Title: title, Image: image}

		jukeBox := &(jockey.JukeBox)
		lock := <-jukeBox.mutex
		jukeBox.songs = append(jukeBox.songs, song)
		jukeBox.mutex <- lock

		fmt.Println("Queue size - %d", len(jukeBox.songs))

		if jukeBox.ended {
			jockey.syncer <- 1
		}

		/*
		Form a JSON here, instead of sending a comma separated string.
		 */
		response, err := json.Marshal(Message{Kind: "Title", Value: id + "," + title + "," + image})
		if err != nil {
			fmt.Println("Unable to form response JSON.")
			response = []byte("Surprise Song!!")
		}

		for _, client := range jockey.Clients {
			client.endpoint.WriteMessage(websocket.TextMessage, response)
		}
	}
}

func (jockey *Jockey) Synchronize() {
	for {
		<-jockey.syncer

		jukeBox := &(jockey.JukeBox)
		lock := <-jukeBox.mutex
		msg := jukeBox.getNextSong()
		jukeBox.mutex <- lock

		time.Sleep(5 * time.Second)

		fmt.Println("Sending next song to all clients %s", string(msg[:]))
		for _, client := range jockey.Clients {
			client.endpoint.WriteMessage(websocket.TextMessage, []byte(msg))
		}

		lock = <-jukeBox.mutex
		jukeBox.startTime = time.Now()
		jukeBox.mutex <- lock
	}
}

func (jukeBox *JukeBox) getNextSong() []byte {
	if (jukeBox.position + 1) == len(jukeBox.songs) {
		jukeBox.ended = true
		msg, _ := json.Marshal(Message{Kind: "Finished"})
		return msg
	} else {
		jukeBox.position ++
		jukeBox.ended = false
		msg, _ := json.Marshal(Message{Kind: "Song", Value: jukeBox.songs[jukeBox.position].Id})
		return msg
	}
}

func getIdAndTitle(query string) (string, string, string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: ""},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Println("Error creating new YouTube client: %v", err)
		return "", "", ""
	}

	response, err := service.Search.List("snippet").MaxResults(1).Q(query).Do()
	if err != nil {
		fmt.Println("Error creating new YouTube client: %v", err)
		return "", "", ""
	}

	return response.Items[0].Id.VideoId, response.Items[0].Snippet.Title, response.Items[0].Snippet.Thumbnails.Default.Url
}
