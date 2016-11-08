package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/floriangrundig/scofw/publisher"
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
	Port                       int
	Hub                        *ws.Hub
	FileEventReportingChannels []chan *publisher.ServerMessage
}

func New(port int, channels []chan *publisher.ServerMessage) *Server {
	return &Server{
		Port: port,
		FileEventReportingChannels: channels,
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

func (server *Server) pushFileChanges() {
	for _, channel := range server.FileEventReportingChannels {
		go func() {
			for {
				// wait for incoming messages
				msg, ok := <-channel

				if !ok {
					log.Println("Shutting down server channel")
					break
				}

				jMsg, err := json.Marshal(&msg)

				if err == nil {
					if *msg.FileChanges != "" {
						server.Hub.Broadcast <- jMsg
					}
				} else {
					log.Println("Error while marshalling event message:", msg)
				}
			}
		}()
	}
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
	go server.pushFileChanges()
	http.ListenAndServe(addr, nil)
}
