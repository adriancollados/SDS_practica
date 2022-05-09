package signs

import (
	"net/http"
	m "sds/util"
)

func Signup(w http.ResponseWriter, req *http.Request) {
	_, ok := m.Gusers[req.Form.Get("user")]

	if ok {
		m.Response(w, false, "Ya existe el usuario", nil)
		return
	}

	u := m.User{}
	u.Name = req.Form.Get("user")

	m.Gusers[u.Name] = u
	m.Response(w, true, "Usuario Registrado", nil)

}
