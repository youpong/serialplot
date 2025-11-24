# serialplot
This sample project illustrates real-time graphing in a browser of values read
from the BBC micro:bit's serial port.
It represents the minimal setup required for visualising sensor values and
similar data in real time.

## Data flow
```
micro:bit -> USB Serial -> PC -> Go Server -> web client
```

## Physical Connection
```
micro:bit - (USB cable) -> PC
```

## Program Roles

### mb/main.py (MicroPython / BBC micro:bit)
- Acquires sensor values on the micro:bit
- **Outputs to USB Serial** as text in the format `A0:number`
- `print()` output reaches the PC's serial port

### main.go(Go / PC)
- Reads data line by line from the serial port
- **Broadcasts in real-time** to all connected WebSocket clients
- Serves `static/index.html` as a web server

### static/index.html(Web client)
- Connects to WebSocket and receives data from Go server
- Renders real-time graph using Chart.js
- Browser operates simply by accessing `http://localhost:8080/`

## File Structure
```
serialplot/
  ├── main.go               # Go: WebSocket server
  ├── static/
  │    └── index.html       # HTML(Chart.js): Web client
  └── mb/
       └── main.py          # MicroPython: micro:bit  
```

## Prerequisites

### Required items
- BBC micro:bit (v1 or v2)
- USB cable
- PC(Go 1.24+)

## Execution Steps

### 1. Upload MicroPython to micro:bit
Copy `mb/main.py` to your micro:bit

### 2. Start the web Server on PC
Get Go modules
```sh
$ go mod tidy
```
Start web server
```sh
$ go run main.go
```

### 3. View in a browser

```
http://localhost:8080/
```
Accessing this page will:
* connect to WebSocket
* real-time drwing using Chart.js

## Data Format
Data sent from micro:bit is **line-based text**

Example:
```
A0:927
A0:1124
```

# LICENSE
MIT
