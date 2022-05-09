package signs

import (
	"crypto"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	util "sds/util"
	"strings"
)

var Claves map[string]crypto.PublicKey

var caracteresInvalidos = map[int]string{0: "!", 1: "\"", 2: "#", 3: "$", 4: "%", 5: "&", 6: "(", 7: ")",
	8: "*", 9: "+", 10: ",", 11: "-", 12: ".", 13: "/", 14: ":", 15: ";", 16: "<", 17: "=",
	18: ">", 19: "?", 20: "@", 21: "[", 22: "\\", 23: "]", 24: "_", 25: "{", 26: "|", 27: "}",
	28: "á", 29: "Á", 30: "é", 31: "É", 32: "í", 33: "Í", 34: "ó", 35: "Ó", 36: "ú", 37: "Ú",
	38: "à", 39: "À", 40: "è", 41: "È", 42: "ì", 43: "Ì", 44: "ò", 45: "Ò", 46: "ù", 47: "Ù",
	48: "ä", 49: "Ä", 50: "ë", 51: "Ë", 52: "ï", 53: "Ï", 54: "ö", 55: "Ö", 56: "ü", 57: "Ü",
	58: "'", 59: "^", 60: "¬", 61: "·"}

func Signup(client *http.Client, cmd string) {
	fmt.Println("Registrar un usuario")
	fmt.Println("--------------------")

	nombreCorrecto := false
	fmt.Print("Nombre de usuario: ")
	user := util.LeerTerminal()

	for !nombreCorrecto {
		nombreCorrecto = true
		for i := 0; i < 62; i++ {
			if strings.Contains(user, caracteresInvalidos[i]) {
				nombreCorrecto = false
				i = 38
				fmt.Println("\nEl nombre de usuario contiene caracteres inválidos")
				fmt.Println("Por favor, repita el nombre de usuario: ")
				user = util.LeerTerminal()
			}
		}
	}

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", cmd)
	data.Set("user", user)
	r, err := client.PostForm("https://localhost:10443", data)
	util.Chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	r.Body.Close()             // hay que cerrar el reader del body
	fmt.Println()
}
