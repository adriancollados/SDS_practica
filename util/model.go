package util

import (
	"crypto/tls"
	"net/http"
	"time"
)

//Cliente global
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var Client = &http.Client{Transport: tr}

var Codee []byte

var TokenSesion []byte

type User struct {
	Name  string // nombre de usuario
	Hash  []byte
	Salt  []byte
	Token []byte
	Data  map[string]string
}

type UsersRegistrados struct {
	Key   []byte
	Users map[string]User
}

type Resp struct {
	Ok    bool
	Msg   string
	Token []byte
}

type Comentario struct {
	Message []byte
	Fecha   time.Time
}

type Fichero struct {
	Name        []byte //nombre fichero
	HashUser    []byte //clave hash del usuario
	Content     []byte //contenido del fichero
	Fecha       time.Time
	Comentarios map[int]Comentario
}

type FicherosRegistrados struct {
	Key      []byte
	Ficheros map[string]Fichero
}

var Gusers map[string]User

var GFicheros map[string]Fichero
