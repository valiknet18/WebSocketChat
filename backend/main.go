package main

import (
	"flag"
	"fmt"
	"log"
	"model/room"
	"net/http"
)

var (
	rooms *Rooms
)

func createRoom(rw http.ResponseWriter, req *http.Request) error {
	if req.Method == "POST" {
		room := &Room{
			name: req.Form["name"]
			register: make(chan *User)
			unregister: make(chan *User)
			sendMessage: make(chan *User)
		}	
	} else {
		http.MethodNotAllowed
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.Handle("/room/create", createRoom)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("Server started ")
}
