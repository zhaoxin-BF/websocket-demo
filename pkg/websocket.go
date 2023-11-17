package pkg

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func HandleWebSocket1(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("hello"))
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}

		time.Sleep(5 * time.Second)
	}
}
