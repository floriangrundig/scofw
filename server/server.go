package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/floriangrundig/scofw/ws"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Server struct {
	Port int
	Hub  *ws.Hub
}

func New(port int) *Server {
	return &Server{
		Port: port,
	}
}

func (server *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handle ws connection request")
	webso, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}

		log.Printf("Error: %v", err)
		return
	}

	client := &ws.Client{Hub: server.Hub, Conn: webso, Send: make(chan []byte, 256)}
	client.Hub.Register <- client
	go client.WritePump()
	client.ReadPump()
}

func (server *Server) createDummyData() {
	ticker := time.NewTicker(1 * time.Second)
	counter := 0
	for {

		select {
		case <-ticker.C:
			server.Hub.Broadcast <- []byte(fmt.Sprintf("Counter: %d", counter))
			counter++
			counter = counter % 60
		}
	}
	// TODO produce some messages and send over hub broadcast channel
}

func (server *Server) Start() {
	fs := http.FileServer(http.Dir("ui/dist"))

	hub := ws.NewHub()
	server.Hub = hub

	go hub.Run()

	http.Handle("/", fs)
	http.HandleFunc("/ws", server.serveWs)

	addr := fmt.Sprintf(":%d", server.Port)

	log.Printf("Starting server at %s", addr)
	go server.createDummyData()
	http.ListenAndServe(addr, nil)
}
