package model

import (
	"encoding/json"
)

type Rooms {
	rooms map[string]Room `json:"rooms"`
	users int32
}

func (r *Rooms) joinRoom(room *Room) {
	r.rooms[room.name] = room
}

type Room struct {
	users    map[string]User `json: "users"`
	messages map[User]Message `json: "messages"`
	name string `json: "name"`
	register chan *User
	unregister chan *User
	sendMessage chan *User
}

func (r *Room) joinToRoom(user *User) {
	r.users[user.name] = user
	Rooms[r.name].users++
}

func (r *Room) leaveFromRoom(user *User) {
	delete(r.users, user.nickname)
	Rooms[r.name].users--
}

type Message struct {
	user *User `json: "User"`
	text string `json: "message"`
}

