package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("Server started ")
}
