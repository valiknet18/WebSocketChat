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

//Структура юзера
type User struct {
	Nickname string          `json:"nickname"`
	Ws       *websocket.Conn `json:"-"`
	RoomHash string          `json:"roomHash"`
	Send     chan []byte     `json:"-"`
	UserHash string          `json:"-"`
}

var users User[]

// Описуем константы для настройки вебсокетов
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

//Метод читает вебсокеты которые приходят
func (u *User) readPump() {
	// log.Println("In user read pump")

	defer func() {
		rooms[u.RoomHash].Unregister <- u
		u.Ws.Close()
		delete(rooms[u.RoomHash].Users, u.UserHash)
		delete(users, u.UserHash)
	}()

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

		t := time.Now()
		currentTime := t.Format("02-01-2006 15:04:05")

		msg := SendMessage{User: users[js.UserHash], Message: js.Message, Time: currentTime}
		message, err = json.Marshal(msg)

		rooms[u.RoomHash].Broadcast <- message
	}
}

// Метод записует данные в вебсокет
func (u *User) write(mt int, payload []byte) error {
	u.Ws.SetWriteDeadline(time.Now().Add(writeWait))

	return u.Ws.WriteMessage(mt, payload)
}

// Метод отлавливает по каналах сообщение, и отправляет юзерам
func (u *User) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		u.Ws.Close()
		delete(rooms[u.RoomHash].Users, u.UserHash)
		delete(users, u.UserHash)
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
	// req.ParseForm()

	// userHash := randString(20)
	// hash := []byte(userHash)

	// user := &User{Nickname: req.Form["nickname"][0], Ws: new(websocket.Conn), RoomHash: req.Form["roomHash"][0], Send: make(chan []byte), UserHash: userHash}

	// // log.Println("Room hash" + string(req.Form["roomHash"][0]))

	// users[userHash] = user

	// log.Println(rooms[req.Form["roomHash"][0]].Users[userHash])
	// rooms[req.Form["roomHash"][0]].Users[userHash] = user

	// rooms[req.Form["roomHash"][0]].Register <- user

	// rw.Header().Set("Content-type", "plain/text")
	// rw.Write(hash)
}

func ConnectToRoom(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ws, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		return
	}

	// fmt.Fprint(rw, req.Form)

	userHash := params.ByName("userHash")

	user := users[userHash]

	user.Ws = ws

	go user.writePump()

	user.readPump()
}

//Функция создает юзера
func UserCreate(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {


	append(users, user)
}