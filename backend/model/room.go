package model

import (
	"encoding/json"
)

type Rooms struct {
	rooms map[string]]Room `json:"rooms"`
	users int32
}

func (r *Rooms) joinRoom(room *Room, hash string) {
	r.rooms[room.name] = room
}

type Room struct {
	users    map[string]User `json: "users"`
	messages map[User]Message `json: "messages"`
	name string `json: "name"`
	register chan *User
	unregister chan *User
	broadcast chan *Message
	roomHash string
}

func (r *Room) joinToRoom(user *User) {
	r.users[user.name] = user
	Rooms[r.name].users++
}

func (r *Room) leaveFromRoom(user *User) {
	delete(r.users, user.nickname)
	Rooms[r.name].users--
}

func (r *Room) run() {
	for {
		select {
			case u: <- r.register:	
				users[u.nickname] = user
			case u: <- r.unregister:
				delete(users, u.nickname)
			case m: <- r.broadcast:
				for _, user := range r.users {
					select {
						case user.message <- m: 
						default: 
							close(user.message)
							delete(users, user.nickname)
					}
				}
		}
	}
}

func connectToRoom(rw http.ResponseWriter, req *http.Request) {
		
}