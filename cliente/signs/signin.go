package signs

import (
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	util "sds/util"
)

func Signin(client *http.Client, cmd string) {
	fmt.Println("Loggear un usuario")
	fmt.Println("--------------------")

	fmt.Print("Nombre de usuario: ")
	user := util.LeerTerminal()

	fmt.Print("Contraseña: ")
	pass := util.LeerTerminal()

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32] // una mitad para el login (256 bits)
	//keyData := keyClient[32:64] // la otra para los datos (256 bits)

	data := url.Values{}                      // estructura para contener los valores
	data.Set("cmd", cmd)                      // comando (string)
	data.Set("user", user)                    // usuario (string)
	data.Set("pass", util.Encode64(keyLogin)) // "contraseña" a base64
	r, err := client.PostForm("https://localhost:10443", data)
	util.Chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	r.Body.Close()             // hay que cerrar el reader del body
	fmt.Println()
}
