package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/websocket"
)

func writeLoop(ws *websocket.Conn, c chan<- int) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print("> ")
		s := scanner.Text()
		if s == "q" {
			c <- 1
			close(c)
			break
		}
		ws.Write([]byte(s))
	}
}

func readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {

		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]
		fmt.Printf("\r%s\n> ", msg)
	}
}

func main() {
	ws, err := websocket.Dial("ws://localhost:3000/ws", "", "http://localhost/")
	c := make(chan int)

	if err != nil {
		fmt.Println("Error dialing server")
	}

	go readLoop(ws)
	go writeLoop(ws, c)

	_, ok := <-c
	if !ok {
		ws.Close()
	}
}
