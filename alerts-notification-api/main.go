package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpointBuilder(fanout *Fanout) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// upgrade this connection to a WebSocket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		log.Println("Client Connected")
		fanout.AddConnection(ws)
		defer fanout.RemoveConnection(ws)

		for {
			// wait until connection is closed
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func setupRoutes(fanout *Fanout) {
	log.Println("Setting up routes")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpointBuilder(fanout))
}

func main() {
	log.Println("Starting alerts notification service")
	ar := NewAlertsReader()
	fanout := NewFanout(&ar.Alerts)
	go fanout.Listen()
	defer fanout.Close()
	go ar.Listen()
	defer ar.Close()

	setupRoutes(fanout)
	log.Fatal(http.ListenAndServe(":12000", nil))
}
