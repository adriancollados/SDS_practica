package models

import (
	"time"
)

type Usuario struct {
	Nombre      string
	Username    string
	Email       string
	Hash        []byte
	Sal         []byte
	Token       []byte
	Ultconexion time.Time
}

type Resp struct {
	Correcto bool
	Mensaje  string
	Token    []byte
}
