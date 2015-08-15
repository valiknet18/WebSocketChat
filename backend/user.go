package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type User struct {
	nickname string
	ws       *websocket.Conn
	message  chan string
}

func (u *User) readPumb() {

}

func sendMessage(rw http.ResponseWriter, req *http.Request) {

}

func connectUser(rw http.ResponseWriter, req *http.Request) {

}
