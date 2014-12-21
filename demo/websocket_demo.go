package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const portString = ":1337"

var connLock sync.Mutex = sync.Mutex{}
var connList []*websocket.Conn = []*websocket.Conn{}

func main() {
	log.Println("Go to http://localhost" + portString)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(portString, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	wsURL := "ws://localhost" + portString + "/ws"
	index := `<!doctype html>
<html>
	<head>
		<title>Websocket Demo</title>
		<script type="text/javascript">
		function loadHandler() {
			var ws = new window.WebSocket("` + wsURL + `");
			ws.onmessage = function(event) {
				var msg = JSON.parse(event.data);
				var text = msg['count'] + '';
				document.getElementById('count').innerHTML = text;
			};
		}
		</script>
	</head>
	<body onload="loadHandler()">
		<label id="count">0</label> clients connected.
	</body>
</html>`
	w.Write([]byte(index))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Got websocket connection.")
	upgrade := websocket.Upgrader{time.Minute, 0, 0, nil, nil, nil}
	conn, err := upgrade.Upgrade(w, r, http.Header{})
	if err != nil {
		log.Print("WebSocket upgrade failed:", err)
		return
	}

	// Broadcast the new connection to everyone.
	connLock.Lock()
	connList = append(connList, conn)
	obj := map[string]int{"count": len(connList)}
	for _, c := range connList {
		c.WriteJSON(obj)
	}
	connLock.Unlock()

	// Wait for the connection to close.
	for {
		// No messages will ever come in, but hey why not.
		var msg interface{}
		if conn.ReadJSON(&msg) != nil {
			break
		}
	}

	// Remove the connection from the list and broadcast the change.
	connLock.Lock()
	obj = map[string]int{"count": len(connList) - 1}
	for i, c := range connList {
		if c == conn {
			connList[i] = connList[len(connList)-1]
			connList = connList[0 : len(connList)-1]
			i--
		} else {
			c.WriteJSON(obj)
		}
	}
	log.Print("Websocket closed. There are ", len(connList), " clients.")
	connLock.Unlock()
}
