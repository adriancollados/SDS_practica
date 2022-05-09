package signs

import (
	"fmt"
	"net/http"
	m "sds/util"
)

func signup(w http.ResponseWriter, req *http.Request) {
	u := m.User{}

	u.Name = req.Form.Get("user")

	_, ok := m.Gusers[u.Name]

	if ok {

		m.Response(w, false, "Ya existe el usuario")
	} else {
		m.Gusers[u.Name] = u
		fmt.Println("Guardado")
		m.Response(w, true, "Usuario registrado con éxito")
	}
}
