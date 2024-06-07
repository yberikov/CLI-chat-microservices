package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to websocket server:", err)
	}

	fmt.Println("starting application")
	done := make(chan os.Signal, 1)
	go chatInputHandler(conn, done)
	go userOutput(conn)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	// Signal received, perform cleanup
	log.Println("Received interrupt signal. Closing all pending connections.")
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "CloseText"))
	if err != nil {
		log.Println("Error during closing websocket:", err)
	}
	err = conn.Close()
	if err != nil {
		log.Println("Error during closing websocket:", err)
	}
}

func userInput(inputChan chan string) {
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	inputChan <- msg
}

func chatInputHandler(conn *websocket.Conn, done chan os.Signal) {
	for {
		inputChan := make(chan string)
		go userInput(inputChan)
		select {
		case input := <-inputChan:
			if strings.TrimSpace(input) == "" {
				continue
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("Error during writing to websocket:", err)
				return
			}

		case <-done:
			log.Println("Interrupt signal received. Closing all pending connections")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Close connection"))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}

func userOutput(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("test")
		}
		fmt.Println(string(message))
	}
}
