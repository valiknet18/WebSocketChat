package model

type User struct {
	nickname string
	ws       *websocket.Conn
}

func (u *User) readPumb() {

}
