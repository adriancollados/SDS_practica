package signs

import (
	"bytes"
	"crypto/rand"
	"net/http"
	m "sds/util"

	"golang.org/x/crypto/argon2"
)

func Signin(w http.ResponseWriter, req *http.Request) {
	u, ok := m.Gusers[req.Form.Get("user")] // ¿existe ya el usuario?
	if !ok {
		m.Response(w, false, "\nERROR: Este usuario no existe en la base de datos", nil)
		return
	}

	password := m.Decode64(req.Form.Get("pass"))
	Hash := argon2.IDKey([]byte(password), u.Salt, 1, 64*1024, 4, 32)

	if !bytes.Equal(u.Hash, Hash) { // comparamos
		m.Response(w, false, "\nERROR: Credenciales inválidas", nil)

	} else {
		u.Token = make([]byte, 16) // token (16 bytes == 128 bits)
		rand.Read(u.Token)         // el token es aleatorio
		m.Gusers[u.Name] = u
		m.Response(w, true, "\n¡Inicio de sesión correcto!", u.Token)
	}
}
