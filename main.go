package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/valiknet18/WebSocketChat/model"
	"log"
	"net/http"
	"text/template"	
	"github.com/rs/cors"
)

func serveHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("static/index.html")

	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, nil)
}

func main() {
	//TODO user httproute package
	r := httprouter.New()

	r.GET("/ws/:userHash/connect", model.ConnectToRoom)
	r.GET("/", serveHome)
	r.POST("/room/create", model.CreateRoom)
	r.GET("/room/get", model.GetRooms)
	r.GET("/room/users/:roomHash", model.GetRoomUsers)
	r.POST("/user/connect", model.ConnectUser)

	r.ServeFiles("/static/*filepath", http.Dir("./static/"))

	flag.Parse()

	handler := cors.Default().Handler(r)

	log.Println("Server running on port: 8000")
	err := http.ListenAndServe(":8001", handler)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
