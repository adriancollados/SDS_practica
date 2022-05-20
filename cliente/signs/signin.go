package signs

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	m "sds/cliente/menu"
	util "sds/util"
)

type UserLogged struct {
	Name string
	Key  []byte
}

var UserLog = UserLogged{}

func Signin(client *http.Client, cmd string) {
	fmt.Println("Loggear un usuario")
	fmt.Println("--------------------")

	fmt.Print("Nombre de usuario: ")
	user := util.LeerTerminal()

	fmt.Print("Contraseña: ")
	pass := util.LeerTerminal()

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	UserLog.Name = user
	UserLog.Key = keyData

	data := url.Values{}                      // estructura para contener los valores
	data.Set("cmd", cmd)                      // comando (string)
	data.Set("user", user)                    // usuario (string)
	data.Set("pass", util.Encode64(keyLogin)) // "contraseña" a base64

	r, err := client.PostForm("https://localhost:10443", data)
	util.Chk(err)

	resp := util.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	fmt.Println(resp.Msg)
	util.TokenSesion = resp.Token
	m.Opciones(resp)
}
