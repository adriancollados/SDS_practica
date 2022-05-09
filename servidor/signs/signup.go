package signs

import (
	"crypto/rand"
	"net/http"
	m "sds/util"

	"golang.org/x/crypto/argon2"
)

func Signup(w http.ResponseWriter, req *http.Request) {
	_, ok := m.Gusers[req.Form.Get("user")]

	if ok {
		m.Response(w, false, "Ya existe el usuario", nil)
		return
	}

	u := m.User{}
	u.Name = req.Form.Get("user")

	u.Salt = make([]byte, 16)
	rand.Read(u.Salt)

	u.Data = make(map[string]string)
	u.Data["private"] = req.Form.Get("prikey")
	u.Data["public"] = req.Form.Get("pubkey")

	password := m.Decode64(req.Form.Get("pass"))

	u.Hash = argon2.IDKey([]byte(password), u.Salt, 1, 64*1024, 4, 32)

	u.Token = make([]byte, 16) // token (16 bytes == 128 bits)
	rand.Read(u.Token)

	m.Gusers[u.Name] = u
	m.Response(w, true, "Usuario Registrado", u.Token)

}
