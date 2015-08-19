package main

import (
	"flag"
	"github.com/valiknet18/WebSocketChat/model"
	"log"
	"net/http"
	"text/template"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
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
	flag.Parse()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/room/create", model.CreateRoom)
	http.HandleFunc("/room/", model.ConnectToRoom)
	http.HandleFunc("/room/get", model.GetRooms)
	// http.Handle("/ws/:room", model.sendMessage)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("Server started ")
}
