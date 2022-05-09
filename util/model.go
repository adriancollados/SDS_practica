package util

import (
	"crypto/tls"
	"net/http"
)

//Cliente global
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var Client = &http.Client{Transport: tr}

type User struct {
	Name  string // nombre de usuario
	Hash  []byte
	Salt  []byte
	Token []byte
	Data  map[string]string
}

type Resp struct {
	Ok    bool
	Msg   string
	Token []byte
}

var Gusers map[string]User
