package util

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
