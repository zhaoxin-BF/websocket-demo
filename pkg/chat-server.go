package pkg

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = make(map[string]*websocket.Conn)
	teams    = make(map[string][]string)
)

func HandleWebSocket2(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	// step1： 登录保持链接
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read message:", err)
	}

	// 解析消息
	perteam, username, _, _, err := parseMessage(msg)
	if err != nil {
		log.Println("Failed to parse message:", err)
	}

	// 记录链接
	clients[username] = conn
	// 记录群聊链接
	teams[perteam] = append(teams[perteam], username)

	// step2: 告知登录成功
	message := "login success"
	err = clients[username].WriteMessage(websocket.TextMessage, []byte(message))

	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		// 解析消息
		perteam, username, peruser, message, err := parseMessage(msg)
		if err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}

		if perteam != "" {
			// 1、保存群聊链接
			//teams[perteam] = append(teams[perteam], username)

			// 2、向群聊发送消息
			for _, user := range teams[perteam] {
				if username != user {
					err = clients[user].WriteMessage(websocket.TextMessage, []byte(message))
					if err != nil {
						log.Println("Failed to send message:", err)
						delete(clients, peruser)
						clients[peruser].Close()
					}
				}
			}
		} else {
			// 2、私聊
			err = clients[peruser].WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Failed to send message:", err)
				delete(clients, peruser)
				clients[peruser].Close()
			} else {
				message = "send success"
				err = clients[username].WriteMessage(websocket.TextMessage, []byte(message))
			}
		}

	}
}

func parseMessage(msg []byte) (string, string, string, string, error) {
	// 解析消息格式，例如 "user=boreas,message=hello"
	// 这里简单地按照逗号进行拆分，您可以根据实际需求进行更复杂的消息解析
	pairs := bytes.Split(msg, []byte(","))
	perteam := ""
	username := ""
	peruser := ""
	message := ""

	for _, pair := range pairs {
		kv := bytes.Split(pair, []byte("="))
		if len(kv) != 2 {
			return "", "", "", "", errors.New("invalid message format")
		}

		key := string(kv[0])
		value := string(kv[1])

		switch key {
		case "perteam":
			perteam = value
		case "username":
			username = value
		case "peruser":
			peruser = value
		case "message":
			message = value
		default:
			return "", "", "", "", errors.New("unknown message key")
		}
	}

	return perteam, username, peruser, message, nil
}
