package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const (
	JoinedServer = "has joined the server"
	LeftServer   = "has left the server"
)

type Server struct {
	conns map[*websocket.Conn]string
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]string),
	}
}

func (s *Server) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)

	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client:", ws.RemoteAddr())

	ws.Write([]byte("Looks like you are a new user, what should we call you?"))
	s.conns[ws] = string(s.readLine(ws))

	s.announce([]byte(s.conns[ws]), JoinedServer)

	s.readLoop(ws)
}

func (s *Server) readLine(ws *websocket.Conn) []byte {
	buf := make([]byte, 1024)
	n, err := ws.Read(buf)
	if err != nil {
		fmt.Println("read error:", err)
	}
	return buf[:n]
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.announce([]byte(s.conns[ws]), LeftServer)
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(s.conns[ws], msg)
	}
}

func (s *Server) announce(user []byte, message string) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write([]byte(fmt.Sprintf("%s %s", user, message))); err != nil {
				fmt.Println(err)
			}
		}(ws)
	}
}

func (s *Server) broadcast(from string, b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write([]byte(fmt.Sprintf("%s: %s", from, b))); err != nil {
				fmt.Println(err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderBook))
	http.ListenAndServe(":3000", nil)
}
