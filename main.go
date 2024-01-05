package main

import (
	"github.com/gorilla/websocket"
	"github.com/zhaoxin-BF/websocket-demo/pkg"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/ws", pkg.HandleWebSocket2)
	http.HandleFunc("/redis")
	log.Println("WebSocket server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
