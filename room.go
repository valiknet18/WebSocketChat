package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

var rooms = make(map[string]*Room)

type Room struct {
	Users      map[string]*User `json:"users"`
	Name       string           `json:"name"`
	Register   chan *User       `json:"-"`
	Unregister chan *User       `json:"-"`
	Broadcast  chan string      `json:"-"`
	RoomHash   string           `json:"hash"`
}

func (r *Room) joinToRoom(user *User) {
	r.Users[user.nickname] = user
}

func (r *Room) leaveFromRoom(user *User) {
	delete(r.Users, user.nickname)
}

func (r *Room) run() {
	for {
		select {
		case u := <-r.Register:
			r.Users[u.nickname] = u
		case u := <-r.Unregister:
			delete(r.Users, u.nickname)
		case m := <-r.Broadcast:
			for _, user := range r.Users {
				select {
				case user.message <- m:
				default:
					close(user.message)
					delete(r.Users, user.nickname)
				}
			}
		}
	}
}

func getRooms(rw http.ResponseWriter, req *http.Request) {
	json, err := json.Marshal(rooms)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-type", "application/json")
	rw.Write(json)
}

func connectToRoom(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
}

//TODO fix this =)
func createRoom(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	if req.Method == "POST" {

		rb := randString(20)

		//TODO fix this name
		room := &Room{
			Name:       req.Form["room"][0],
			Register:   make(chan *User),
			Unregister: make(chan *User),
			Broadcast:  make(chan string),
			RoomHash:   rb,
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
