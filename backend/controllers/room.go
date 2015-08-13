package room

import (
	"log"
)

type Connection struct {
}

type Room struct {
	users    []User
	messages []Message
}

type User struct {
	nickname string
}

type Message struct {
	user *User
	text string
}

func (r *room) joinToRoom(u Users) {
	append(r.users, u)
}
