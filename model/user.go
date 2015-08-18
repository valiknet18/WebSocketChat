package model

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type User struct {
	nickname string
	ws       *websocket.Conn
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
	defer func() {
		// .unregister <- c
		u.ws.Close()
	}()

	u.ws.SetReadLimit(maxMessageSize)
	u.ws.SetReadDeadline(time.Now().Add(pongWait))
	u.ws.SetPongHandler(func(string) error { u.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// for {
	// 	_, message, err := u.ws.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	// h.broadcast <- message
	// }
}

// write writes a message with the given message type and payload.
func (u *User) write(mt int, payload []byte) error {
	u.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return u.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (u *User) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.ws.Close()
	}()
	// for {
	// 	select {
	// 	case message, ok := <-u.send:
	// 		if !ok {
	// 			u.write(websocket.CloseMessage, []byte{})
	// 			return
	// 		}
	// 		if err := u.write(websocket.TextMessage, message); err != nil {
	// 			return
	// 		}
	// 		// case <-ticker.C:
	// 		// 	if err := u.write(websocket.PingMessage, []byte{}); err != nil {
	// 		// 		return
	// 		// 	}
	// 	}
	// }
}

func SendMessage(rw http.ResponseWriter, req *http.Request) {

}

func SonnectUser(rw http.ResponseWriter, req *http.Request) {

}
