package radio

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"bytes"
	//"regexp"
)

const MaximumClients = 5

type Client struct {
	Handle   string
	Endpoint *websocket.Conn
}

func (client Client) Listen(queue chan string) {
	fmt.Println("Running Listener ...")
	for {
		messageType, data, err := client.Endpoint.ReadMessage()
		if err != nil {
			fmt.Println("Error on reading socket. Closing the listener.")
			break
		}
		fmt.Println(data)

		if messageType == websocket.TextMessage {
			inputString := string(data[:])
			fmt.Println("Received input - " + inputString)

			//m, _ := regexp.MatchString("^https://youtube\\.com/", inputString)
			//if m {
			queue <- inputString
			//}
		}
	}
}

type ClientList struct {
	Clients      [MaximumClients]Client
	NumOfClients int
	Queue        chan string
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
	go newClient.Listen(clients.Queue)
	clients.Clients[clients.NumOfClients] = Client{Endpoint: conn}
	clients.NumOfClients ++
	fmt.Println(clients.NumOfClients)
}

func (clients *ClientList) PopulateQueue(songList *[]string) {
	fmt.Println("Reading channel ...")
	for {
		v := <-clients.Queue
		fmt.Println("One of the clients entered - " + v)

		for _, client := range clients.Clients[:clients.NumOfClients] {
			client.Endpoint.WriteMessage(websocket.TextMessage, []byte(v))
		}
	}
}
