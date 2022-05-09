package util

type User struct {
	Name string // nombre de usuario
}

type Resp struct {
	Ok  bool
	Msg string
}

var Gusers map[string]User
