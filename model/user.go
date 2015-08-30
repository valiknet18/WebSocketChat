package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type User struct {
	Nickname string          `json:"nickname"`
	Ws       *websocket.Conn `json:"-"`
	RoomHash string          `json:"roomHash"`
	Send     chan []byte     `json:"-"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (u *User) readPump() {
	// log.Println("In user read pump")

	defer func() {
		// u.unregister <- c
		u.Ws.Close()
	}()

	// defer u.Ws.Close()

	fmt.Println(u.Ws)

	u.Ws.SetReadLimit(maxMessageSize)
	u.Ws.SetReadDeadline(time.Now().Add(pongWait))
	u.Ws.SetPongHandler(func(string) error { u.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := u.Ws.ReadMessage()

		if err != nil {
			break
		}

		js := new(ReturnMessage)

		err = json.Unmarshal(message, &js)

		msg := SendMessage{User: users[js.UserHash], Message: js.Message}

		message, err = json.Marshal(msg)

		rooms[u.RoomHash].Broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (u *User) write(mt int, payload []byte) error {
	// fmt.Println(string(payload))

	u.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return u.Ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (u *User) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		u.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-u.Send:
			{
				if !ok {
					u.write(websocket.CloseMessage, []byte{})
				}

				if err := u.write(websocket.TextMessage, message); err != nil {
					return
				}
			}

		case <-ticker.C:
			{
				if err := u.write(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}
}

// func SendMessage(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
// }

func ConnectUser(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()

	userHash := randString(20)
	hash := []byte(userHash)

	user := &User{Nickname: req.Form["nickname"][0], Ws: new(websocket.Conn), RoomHash: req.Form["roomHash"][0], Send: make(chan []byte)}

	// log.Println("Room hash" + string(req.Form["roomHash"][0]))

	users[userHash] = user

	log.Println(rooms[req.Form["roomHash"][0]].Users[userHash])
	rooms[req.Form["roomHash"][0]].Users[userHash] = user

	rw.Header().Set("Content-type", "plain/text")
	rw.Write(hash)
}
