package util

type User struct {
	Name string // nombre de usuario
}

type Resp struct {
	Ok    bool
	Msg   string
	Token []byte
}

var Gusers map[string]User
