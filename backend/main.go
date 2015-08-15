package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
)

//Create new room and run him in goroutines
func createRoom(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

		rb := make([]byte, 32)
		_, err := rand.Read(rb)

		if err != nil {
			fmt.Println(err)
		}

		rs := base64.URLEncoding.EncodeToString(rb)

		room := &Room{
			name:       req.Form["name"],
			register:   make(chan *User),
			unregister: make(chan *User),
			broadcast:  make(chan string),
			roomHash:   rs,
		}

		go room.run()

		rooms.joinRoom(room, rs)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/user/connect", connectUser)
	http.HandleFunc("/room/create", createRoom)
	http.HandleFunc("/room/connect", connectToRoom)
	// http.Handle("/ws/:room", model.sendMessage)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("Server started ")
}
