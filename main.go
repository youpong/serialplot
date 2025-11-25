package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

type MockReader struct{}

func (m *MockReader) Read(p []byte) (int, error) {
	s := fmt.Sprintf("%d,%d\n", rand.Intn(1000), rand.Intn(1000))
	time.Sleep(100 * time.Millisecond)
	// fmt.Printf("DEBUG: %s", s)
	return copy(p, s), nil
}

func main() {
	var port_name string
	var mock bool
	flag.StringVar(&port_name, "port", "/dev/tty.usbmodem1101", "serial port")
	flag.BoolVar(&mock, "mock", false, "use mock data instead of serial port(develop)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\nExample serial port name:\n"+
			"\tLinux:   /dev/ttyACM0\n"+
			"\tmacOS:   /dev/tty.usbmodem1101\n"+
			"\tWindows: COM3\n")
	}
	flag.Parse()

	var src io.Reader
	if !mock {
		c := &serial.Config{
			Name: port_name,
			Baud: 115200, // micro:bit standard baud
		}
		src, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
	} else {
		src = &MockReader{}
	}
	scanner := bufio.NewScanner(src) // TODO: scanner の後始末って必要？

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
		var accel, temp int
		fmt.Sscanf(line, "%d,%d", &accel, &temp)
		//		fmt.Printf("DEBUG: %d, %d\n", accel, temp)
		payload := map[string]interface{}{
			"values": map[string]int{
				"accel": accel,
				"temp":  temp,
			},
		}

		jsonBytes, _ := json.Marshal(payload)
		broadcast(jsonBytes)
	}
}
