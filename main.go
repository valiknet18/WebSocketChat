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
	if r.URL.Path == "/static/js/app.js" {
		http.ServeFile(w, r, "static/js/app.js")

		return
	}

	if r.URL.Path == "/static/css/app.css" {
		http.ServeFile(w, r, "static/css/app.css")

		return
	}

	if r.URL.Path == "/static/templates/chat.html" {
		http.ServeFile(w, r, "static/templates/chat.html")

		return
	}

	homeTempl := template.Must(template.ParseFiles("static/index.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
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
	err := http.ListenAndServe(":8000", handler)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
