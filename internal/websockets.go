package internal

import (
	"bytes"
	"github.com/gofiber/websocket/v2"
	"html/template"
	"log"
)

type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *Message
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *Message),
	}
}

func (s *WebSocketServer) HandleWebSocket(conn *websocket.Conn) {
	s.clients[conn] = true
	defer func() {
		delete(s.clients, conn)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Printf("error occurred while writing to a client: %v ", err)
				delete(s.clients, client)
				client.Close()
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		log.Fatalf("error occurred while parsing template: %s", err)
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("error occurred while executing template: %s", err)
	}

	return renderedMessage.Bytes()
}
