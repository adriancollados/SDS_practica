package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	u "sds/util"
	"strings"
	"time"
)

func main() {
	fmt.Println("---------------------------------------")
	fmt.Println("Bienvenido al sistema de SDS (CLIENTE)")
	fmt.Println("---------------------------------------")

	for {
		fmt.Println("\nIntroduzca una opción:")
		fmt.Println("-----------------------")
		fmt.Println("1. Registrar Usuario")
		fmt.Println("2. Iniciar Sesión")
		fmt.Println("-----------------------")
		fmt.Println("3. Salir del programa")
		fmt.Print("\nOpción: ")
		option := u.LeerTerminal()
		if option == "1" || option == "2" || option == "3" {
			switch option {
			case "1":
				Signup(u.Client, "signup")
			case "2":
				Signin(u.Client, "signin")
			case "3":
				u.TokenSesion = nil
				fmt.Println("\n¡Hasta luego!")
				return
			}
		} else {
			fmt.Println("No es una opción válida introduzca un número entre 1 y 3:")
		}

		for {
			fmt.Print("\n¿Desea realizar otra operación? (s/n): ")
			continuar := u.LeerTerminal()

			if continuar != "s" && continuar != "n" {
				fmt.Println("\nPor favor, introduzca una respuesta válida")
			} else if continuar == "n" {
				fmt.Println("\n¡Hasta luego!")
				return
			} else if continuar == "s" {
				break
			}
		}

	}
}

///////////////////////////////////////////
///////				SINGUP			///////
///////////////////////////////////////////

var Claves map[string]crypto.PublicKey

//esto se puede cambiar por una funcion del paquete de string
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
	fmt.Print("Introduzca un nombre de usuario: ")
	user := u.LeerTerminal()

	for !nombreCorrecto {
		nombreCorrecto = true
		for i := 0; i < 62; i++ {
			if strings.Contains(user, caracteresInvalidos[i]) {
				nombreCorrecto = false
				i = 38
				fmt.Println("\nEl nombre de usuario contiene caracteres inválidos")
				fmt.Println("Por favor, repita el nombre de usuario: ")
				user = u.LeerTerminal()
			}
		}
	}

	fmt.Print("Introduzca su contraseña: ")
	pass := u.LeerTerminal()
	fmt.Println()

	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32]
	keyData := keyClient[32:64] //La otra para los datos

	//Generamos un par de claves (privada, pública) para el servidor
	pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
	u.Chk(err)
	pkClient.Precompute()

	pkJSON, err := json.Marshal(&pkClient)
	u.Chk(err)

	keyPub := pkClient.Public()
	pubJSON, err := json.Marshal(&keyPub)
	u.Chk(err)

	// println(keyPub)
	// println(pubJSON)
	// print(string(pubJSON))

	//Guardamos en un fichero la clave publica
	//os.Create("cp.json")
	//se puede cambiar por os.WriteFile
	err = ioutil.WriteFile("cp.json", pubJSON, 0666)
	if err != nil {
		fmt.Println(err)
	}

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("user", user)
	data.Set("pass", u.Encode64(keyLogin))

	//Comprimimos y codificamos la clave pública
	data.Set("pubkey", u.Encode64(u.Compress(pubJSON)))

	//Comprimimos ciframos y codificamos la clave privada
	data.Set("prikey", u.Encode64(u.Encrypt(u.Compress(pkJSON), keyData)))

	r, err := client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	if resp.Ok {
		fmt.Println(resp.Msg)
	}
}

///////////////////////////////////////////
///////				SIGNIN			///////
///////////////////////////////////////////

type UserLogged struct {
	Name string
	Key  []byte
}

var UserLog = UserLogged{}

func Signin(client *http.Client, cmd string) {
	fmt.Println("Loggear un usuario")
	fmt.Println("--------------------")

	fmt.Print("Nombre de usuario: ")
	user := u.LeerTerminal()

	fmt.Print("Contraseña: ")
	pass := u.LeerTerminal()

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	UserLog.Name = user
	UserLog.Key = keyData

	data := url.Values{}                   // estructura para contener los valores
	data.Set("cmd", cmd)                   // comando (string)
	data.Set("user", user)                 // usuario (string)
	data.Set("pass", u.Encode64(keyLogin)) // "contraseña" a base64

	r, err := client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)

	if resp.Ok {
		fmt.Println(resp.Msg)
		u.TokenSesion = resp.Token
	}
	Opciones(resp)
}

///////////////////////////////////////////
///////			CrearFich			///////
///////////////////////////////////////////

func CrearFich(cmd string) {
	fmt.Println("Creando fichero ...")
	fmt.Println("------------------------")

	if u.GFicheros == nil {
		u.GFicheros = make(map[string]u.Fichero)
	}

	f := u.Fichero{}
	fmt.Println("Nombre del fichero: ")
	f.Name = u.LeerTerminal()
	fmt.Println("Contenido del fichero: ")
	f.Content = []byte(u.LeerTerminal())
	f.HashUser = UserLog.Key

	_, ok := u.GFicheros[f.Name]
	if ok {
		fmt.Println("El fichero ya existe")
	} else {
		u.GFicheros[f.Name] = f
	}

	jsonData, err := json.Marshal(&f)
	u.Chk(err)
	jsonData = []byte(u.Encode64(u.Encrypt(jsonData, UserLog.Key)))

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("fich", string(jsonData))
	data.Set("user", string(UserLog.Key))

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp)
	// u.GFicheros = make(map[string]u.Fichero)
	fmt.Println("Fichero creado")
	fmt.Println("------------------------")
	Opciones(resp)
}

///////////////////////////////////////////
///////			EnvioFich			///////
///////////////////////////////////////////

func Fichup(filename string, cmd string) {
	filepath := "../archivos/" + filename
	if isExist(filepath) {
		fmt.Println("La ruta es correcta")
		// file, err := os.Open(filepath) // abrimos el fichero de origen (cliente)
		// u.Chk(err)
		// defer file.Close() // cerramos al salir de ámbito

		t := time.Now() // timestamp para medir tiempo
		// hacemos un post con formato octet-stream (binario) y ponemos el reader del fichero directamente como Body
		data := url.Values{}
		data.Set("cmd", cmd) // comando (string)
		data.Set("fichero", filename)
		data.Set("user", u.Encode64(UserLog.Key))

		resp, err := u.Client.PostForm("https://localhost:10443", data)
		u.Chk(err)
		fmt.Println("CLIENTE::", time.Since(t), ":: Post realizado") // imprimimos tiempo
		io.Copy(os.Stdout, resp.Body)                                // copiamos la respuesta a la salida estándar
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

///////////////////////////////////////////
///////			LEERFICH			///////
///////////////////////////////////////////

func LeerFich(cmd string, filename string) {
	fmt.Println("Leyendo fichero ...")
	fmt.Println("------------------------")

	if u.GFicheros == nil {
		u.GFicheros = make(map[string]u.Fichero)
	}

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("filename", filename)
	data.Set("user", u.Encode64(UserLog.Key))

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp)

	var fich = u.Fichero{}
	json.Unmarshal([]byte(resp.Msg), &fich.Content)

	if resp.Ok {
		fmt.Println("Nombre del fichero: " + filename)
		fmt.Println("Contenido: " + resp.Msg)
	} //else {
	// 	fmt.Println(resp.Msg)
	// }

	Opciones(resp)
}

///////////////////////////////////////////
///////				Menu			///////
///////////////////////////////////////////
func Opciones(resp u.Resp) {
	if !resp.Ok {
		fmt.Println(resp.Msg)
		return
	} else {
		for {
			fmt.Println("\n---- MENÚ PRINCIPAL ----")
			fmt.Println("1. Crear archivo")
			fmt.Println("2. Leer archivo")
			fmt.Println("3. Subir archivo")
			fmt.Println("4. Cerrar sesión")
			fmt.Println("------------------------")
			fmt.Print("¿Qué opción desea realizar? ")
			option := u.StringAInt(u.LeerTerminal())
			fmt.Println("")

			switch option {
			case 1:
				fmt.Println("Se ha seleccionado CREAR ARCHIVO")
				fmt.Println("--------------------------------")
				CrearFich("crearFichero")
			case 2:
				fmt.Println("Se ha seleccionado LEER ARCHIVO")
				fmt.Println("--------------------------------")
				fmt.Print("Introduzca el nombre del fichero que desea ver: ")
				filename := u.LeerTerminal()
				LeerFich("leerFichero", filename)
			case 3:
				fmt.Println("Se ha seleccionado SUBIR ARCHIVO")
				fmt.Println("--------------------------------")
				fmt.Print("Introduzca el nombre del fichero que desea subir: ")
				filename := u.LeerTerminal()

				Fichup(filename, "subirFichero")
			case 4:
				fmt.Println("\n¡Hasta luego!")
				return
			default:
				fmt.Println("No es una opción válida introduzca un número entre 1 y 3:")
				Opciones(resp)
				return
			}

		}

	}

}
