package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO: いるの？
}

var wsConnections = make(map[*websocket.Conn]bool)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade error", http.StatusInternalServerError)
		return
	}
	defer conn.Close() // TODO: conn.Close() を呼び出すタイミングを検討しよう。ここが最適か？
	// TODO: ここで呼び出す必要があるのか？
	wsConnections[conn] = true
	for {
		if _, _, err := conn.NextReader(); err != nil {
			// TODO: ここで conn.Close() するのでは？
			delete(wsConnections, conn)
			return
		}
	}
}

func broadcast(msg []byte) {
	for conn := range wsConnections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.Close()
			delete(wsConnections, conn)
		}
	}
}

func main() {
	// Example micro:bit serial port name
	//   Linux:   /dev/ttyACM0"
	//   macOS:   /dev/tty.usbmodem1101
	//   Windows: COM3
	// TODO: これはファイルの先頭に持っていきたい。
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

	scanner := bufio.NewScanner(port) // TODO: scanner の後始末って必要？

	// Start web server
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// TODO: go routine 化する必要ある？
	go func() {
		log.Println("Listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Loop to broadcast values read from the serial port
	for scanner.Scan() {
		line := scanner.Text()
		// Verification format label:value
		var accel, gyro, temp int
		fmt.Sscanf(line, "%d,%d,%d", &accel, &gyro, &temp)

		payload := map[string]interface{}{
			"values": map[string]int{
				"accel": accel,
				"gyro":  gyro,
				"temp":  temp,
			},
		}

		jsonBytes, _ := json.Marshal(payload)
		broadcast(jsonBytes)
	}
}
