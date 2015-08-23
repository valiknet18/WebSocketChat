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
)

var rooms = make(map[string]*Room)
var users = make(map[string]*User)

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

func (r *Room) joinToRoom(user *User) {
	r.Users[user.Nickname] = user
}

func (r *Room) leaveFromRoom(user *User) {
	delete(r.Users, user.Nickname)
}

func (r *Room) run() {
	for {
		select {
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

func GetRooms(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	json, err := json.Marshal(rooms)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-type", "application/json")
	rw.Write(json)
}

func ConnectToRoom(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ws, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		return
	}

	fmt.Fprint(rw, req.Form)

	_, message, _ := ws.ReadMessage()

	log.Println(req)
	log.Println(message)

	userHash := params.ByName("userHash")

	user := users[userHash]

	go user.writePump()

	user.readPump()

	// h.register <- c
	// go c.writePump()
	// c.readPump()
}

func print_binary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

//TODO fix this =)
func CreateRoom(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()

	if req.Method == "POST" {

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

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GetRoomUsers(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	hash := params.ByName("roomHash")

	js, err := json.Marshal(rooms[hash])

	if err != nil {
		log.Println(err.Error)
	}

	rw.Header().Set("Content-type", "application/json")
	rw.Write(js)
}
