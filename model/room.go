//TODO Make struct for json result

package model

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

var rooms = make(map[string]*Room)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Room struct {
	Users      map[string]*User `json:"users"`
	Name       string           `json:"name"`
	Register   chan *User       `json:"-"`
	Unregister chan *User       `json:"-"`
	Broadcast  chan []byte      `json:"-"`
	RoomHash   string           `json:"hash"`
}

//Присойденяет юзера к комнате
func (r *Room) joinToRoom(user *User) {
	r.Users[user.Nickname] = user
}

//Удаляет юзера с комнаты
func (r *Room) leaveFromRoom(user *User) {
	delete(r.Users, user.Nickname)
}

//Запускает комнату
func (r *Room) run() {
	t := time.Now()
	currentTime := t.Format("02-01-2006 15:04:05")

	for {
		select {
		case u := <-r.Register:
			{

				msg := SendMessage{User: u, Message: "Присоединился к комнате", Time: currentTime}
				message, err := json.Marshal(msg)

				if err != nil {
					fmt.Println(err)
				}

				for _, u := range r.Users {
					fmt.Println(u.Nickname)

					select {
					case u.Send <- message:
					}
				}
			}
			// r.Users[u.UserHash] = u

		case u := <-r.Unregister:
			{
				msg := SendMessage{User: u, Message: "Покинул комнату", Time: currentTime}
				message, err := json.Marshal(msg)

				if err != nil {
					fmt.Println(err)
				}

				for _, u := range r.Users {
					fmt.Println(u.Nickname)

					select {

					case u.Send <- message:
					}
				}
			}

		case m := <-r.Broadcast:
			{
				for _, u := range r.Users {
					select {

					case u.Send <- m:

					default:
						{
							close(u.Send)
						}
					}
				}
			}
		}
	}
}

//Возвращает все комнаты
func GetRooms(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	json, err := json.Marshal(rooms)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-type", "application/json")
	rw.Write(json)
}

func print_binary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

//Создает новую комнату
func CreateRoom(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()

	if req.Method == "POST" {

		fmt.Println("Name: ", req.FormValue("room"))

		rb := randString(20)

		//TODO fix this name
		room := &Room{
			Name:       req.Form["room"][0],
			Register:   make(chan *User),
			Unregister: make(chan *User),
			Broadcast:  make(chan []byte),
			RoomHash:   rb,
			Users:      make(map[string]*User),
		}

		go room.run()

		rooms[rb] = room

		js, err := json.Marshal(rooms)

		if err != nil {
			fmt.Fprintf(rw, err.Error(), http.StatusInternalServerError)
		}

		rw.Header().Set("Content-type", "application/json")
		rw.Write(js)
	}
}

//Рандомная строка
func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

//Возвращает юзеров комнаты
func GetRoomUsers(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	hash := params.ByName("roomHash")

	js, err := json.Marshal(rooms[hash])

	if err != nil {
		log.Println(err.Error)
	}

	rw.Header().Set("Content-type", "application/json")
	rw.Write(js)
}
