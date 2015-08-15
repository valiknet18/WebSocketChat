package main

import (
	// "encoding/json"
	"net/http"
)

var (
	rooms *Rooms
)

type Rooms struct {
	rooms map[string]*Room `json:"rooms"`
	users int32
}

func (r *Rooms) joinRoom(room *Room, hash string) {
	r.rooms[hash] = room
}

type Room struct {
	users      map[string]*User `json: "users"`
	messages   map[User]string  `json: "messages"`
	name       []string         `json: "name"`
	register   chan *User
	unregister chan *User
	broadcast  chan string
	roomHash   string
}

func (r *Room) joinToRoom(user *User) {
	r.users[user.nickname] = user
}

func (r *Room) leaveFromRoom(user *User) {
	delete(r.users, user.nickname)
}

func (r *Room) run() {
	for {
		select {
		case u := <-r.register:
			r.users[u.nickname] = u
		case u := <-r.unregister:
			delete(r.users, u.nickname)
		case m := <-r.broadcast:
			for _, user := range r.users {
				select {
				case user.message <- m:
				default:
					close(user.message)
					delete(r.users, user.nickname)
				}
			}
		}
	}
}

func connectToRoom(rw http.ResponseWriter, req *http.Request) {

}
