package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
	"crypto/md5"
	"fmt"
	"encoding/hex"
)

//Структура юзера
type User struct {
	Nickname string          `json:"nickname"`
	Ws       *websocket.Conn `json:"-"`
	RoomHash string          `json:"roomHash"`
	Send     chan []byte     `json:"-"`
	UserHash string          `json:"-"`
}

var users = make(map[string]*User)

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


func ConnectToRoom(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ws, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		return
	}

	roomHash := params.ByName("roomHash")
	userHash := params.ByName("userHash")

	user := users[userHash]
	room := rooms[roomHash]

	room.joinToRoom(user)
	user.Ws = ws

	go user.writePump()

	user.readPump()
}

//Функция создает юзера
func UserCreate(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	req.ParseForm()
	var salt string = "ke12"

	var nickname string = req.FormValue("nickname")

	data := []byte(nickname + salt)
	hash := md5.Sum(data)

	users[string(hash[:])] = &User{Nickname: nickname}

	jsStruct := struct{ Hash string } { Hash: hex.EncodeToString(hash[:]) }

	js, err := json.Marshal(jsStruct)

	if err != nil {
		panic(err)
	}

	fmt.Println("Hash code: ", hex.EncodeToString(hash[:]))

	rw.Header().Set("Content-type", "application/json")
	rw.Write(js)	
}