package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"go.bug.st/serial"
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
	// ---- 1. シリアルポートを開く ----
	port, err := serial.Open(
		"/dev/ttyACM0", // Windows: COM3 等
		&serial.Mode{BaudRate: 115200},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	scanner := bufio.NewScanner(port)

	// ---- 2. Web サーバ起動 ----
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	go func() {
		log.Println("Web server started on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// ---- 3. シリアル読み取りループ ----
	for scanner.Scan() {
		line := scanner.Text()
		// label:value の形式確認
		if strings.Contains(line, ":") {
			broadcast(line)
			fmt.Println("Recv:", line)
		}
	}
}
