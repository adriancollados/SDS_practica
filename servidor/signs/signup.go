package signs

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	m "sds/util"

	"golang.org/x/crypto/argon2"
)

func Signup(w http.ResponseWriter, req *http.Request) {
	_, ok := m.Gusers[req.Form.Get("user")]

	if ok {
		m.Response(w, false, "\nERROR: Este usuario ya existe en la base de datos", nil)
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
	var code []byte = nil
	User := m.UsersRegistrados{Key: code, Users: m.Gusers}
	User.Key = m.Codee
	User.Users = m.Gusers
	os.Remove("users.json")
	_, err := os.Create("users.json")
	m.Chk(err)
	jsonF, err := json.Marshal(&User)
	m.Chk(err)
	//Encriptamos el json de los ficheros con el codigo de la contraseña del server
	var jsonFD = jsonF

	err = ioutil.WriteFile("users.json", jsonFD, 0644)
	m.Chk(err)
	m.Response(w, true, "\n¡Usuario registrado correctamente!", u.Token)

}
