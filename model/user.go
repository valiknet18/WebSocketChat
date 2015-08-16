package model

type User struct {
	nickname string
	ws       *websocket.Conn
	message  chan *Message
}

func (u *User) readPumb() {

}

func sendMessage(rw http.ResponseWriter, req *http.Request) {

}

func connectUser(rw http.ResponseWriter, req *http.Request) {

}
