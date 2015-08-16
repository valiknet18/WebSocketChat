package main

import (
	"flag"
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

	homeTempl := template.Must(template.ParseFiles("static/index.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/user/connect", connectUser)
	http.HandleFunc("/room/create", createRoom)
	http.HandleFunc("/room/connect", connectToRoom)
	http.HandleFunc("/room/get", getRooms)
	// http.Handle("/ws/:room", model.sendMessage)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("Server started ")
}
