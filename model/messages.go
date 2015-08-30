package model

type ReturnMessage struct {
	UserHash string
	Message  string
}

type SendMessage struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}
