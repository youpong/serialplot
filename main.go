package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var wsConnections = make(map[*websocket.Conn]bool)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade error", http.StatusInternalServerError)
		return
	}
	wsConnections[conn] = true
	defer conn.Close()
	for {
		if _, _, err := conn.NextReader(); err != nil {
			delete(wsConnections, conn)
			return
		}
	}
}

func broadcast(msg string) {
	for conn := range wsConnections {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func main() {
	// Example micro:bit serial port name
	//   Linux:   /dev/ttyACM0"
	//   macOS:   /dev/tty.usbmodem1101
	//   Windows: COM3
	c := &serial.Config{
		Name: "/dev/tty.usbmodem1101",
		Baud: 115200, // micro:bit standard baud
	}

	if len(os.Args) == 2 {
		c.Name = os.Args[1]
	}

	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	scanner := bufio.NewScanner(port)

	// Start web server
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	go func() {
		log.Println("Web server started on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Loop to broadcast values read from the serial port
	for scanner.Scan() {
		line := scanner.Text()
		// Verification format label:value
		if strings.Contains(line, ":") {
			broadcast(line)
			fmt.Println("Recv:", line)
		}
	}
}
