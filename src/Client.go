package radio

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	handle   string
	id       int
	endpoint *websocket.Conn
}

func (client Client) Listen(queue chan string, close chan int, syncer chan int) {
	fmt.Println("Running Listener ...")
	for {
		messageType, data, err := client.endpoint.ReadMessage()
		if err != nil {
			fmt.Println("Error on reading socket. Closing the listener.")
			close <- client.id
			break
		}

		if messageType == websocket.TextMessage {
			inputString := string(data[:])
			fmt.Println("One of the clients entered - " + inputString)

			if "Finished" == inputString {
				select {
				case syncer <- 1:
					fmt.Println("Running Synchronizer")
				default:
					fmt.Println("Nothing to do")
				}
			} else {
				if "" != inputString {
					/*
					The condition is needed as clicking Add to Queue with nothing input is causing the
					YouTube API to send a random song. Need to check why.
					 */
					queue <- inputString
				}
			}
		}
	}
}
