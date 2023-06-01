package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Fanout struct {
	wsConnections []*websocket.Conn
	alerts        *chan []byte
	mu            sync.Mutex
	openChan      chan bool
}

func NewFanout(alerts *chan []byte) *Fanout {
	return &Fanout{
		wsConnections: make([]*websocket.Conn, 0),
		alerts:        alerts,
		mu:            sync.Mutex{},
		openChan:      make(chan bool),
	}
}

func (f *Fanout) AddConnection(conn *websocket.Conn) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.wsConnections = append(f.wsConnections, conn)
}

func (f *Fanout) RemoveConnection(conn *websocket.Conn) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, c := range f.wsConnections {
		if c == conn {
			f.wsConnections = append(f.wsConnections[:i], f.wsConnections[i+1:]...)
			break
		}
	}
}

func (f *Fanout) Listen() {
	for {
		select {
		case <-f.openChan:
			log.Println("Closing fanout")
			return
		case alert := <-*f.alerts:
			log.Printf("Sending alert %s to all connections\n", string(alert))
			f.mu.Lock()
			defer f.mu.Unlock()
			for _, conn := range f.wsConnections {
				conn.WriteMessage(1, alert)
			}
		}
	}
}

func (f *Fanout) Close() {
	f.openChan <- false
}
