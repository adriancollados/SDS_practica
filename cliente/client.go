package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	u "sds/util"
	"strings"
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
		contains := false
		for i := 0; i < 62 && contains == false; i++ {
			if strings.Contains(user, caracteresInvalidos[i]) || user == "" {
				contains = true
			}
		}

		if contains == true {
			fmt.Println("\nEl nombre de usuario contiene caracteres inválidos")
			fmt.Println("Por favor, repita el nombre de usuario: ")
			user = u.LeerTerminal()
		} else {
			nombreCorrecto = true
		}
	}

	fmt.Println()
	fmt.Print("Introduzca su contraseña: ")
	pass := u.LeerTerminal()
	fmt.Println()

	for pass == "" {
		fmt.Println("\nLa contraseña no puede ser vacía")
		fmt.Println("Por favor, repita la contaseña: ")
		pass = u.LeerTerminal()
	}

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
		UserLog.Name = user
		UserLog.Key = keyData
	}
	Opciones(resp)
}

///////////////////////////////////////////
///////			CrearFich			///////
///////////////////////////////////////////

func CrearFich(cmd string) {
	fmt.Println("Creando fichero ...")
	fmt.Println("------------------------")

	idValid := false
	f := u.Fichero{}
	fmt.Println("Nombre del fichero: ")
	id := u.LeerTerminal()

	for !idValid {
		contains := false
		for i := 0; i < 62 && contains == false; i++ {
			if strings.Contains(id, caracteresInvalidos[i]) || id == "" {
				contains = true
			}
		}

		if contains == true {
			fmt.Println("\nEl nombre de fichero contiene caracteres inválidos")
			fmt.Println("Por favor, repita el nombre de fichero: ")
			id = u.LeerTerminal()
		} else {
			idValid = true
			id = id + ".txt"
		}
	}

	f.Name = []byte(id)
	fmt.Println("Contenido del fichero: ")
	f.Content = []byte(u.LeerTerminal())
	f.HashUser = UserLog.Key

	comentario := false
	coment := u.Comentario{}
	data := url.Values{}
	for !comentario {
		fmt.Println("Desea añadir algún comentario a su fichero? (s/n): ")
		respuesta := u.LeerTerminal()

		if respuesta == "s" {
			comentario = true
			fmt.Println("Introduzca el comentario: ")
			coment.Message = []byte(u.LeerTerminal())

			jsonComment := u.Encode64(u.Encrypt(coment.Message, UserLog.Key))
			data.Set("comments", jsonComment)

		} else if respuesta == "n" {
			comentario = true
		} else {
			fmt.Println("\nPor favor, introduzca una respuesta válida")
		}
	}

	jsonName := u.Encode64(u.Encrypt([]byte(f.Name), UserLog.Key))
	var aux = string(u.Decode64(jsonName))
	var entra = false
	for ok := true; ok; ok = strings.ContainsAny(aux, "/") {
		aux = u.Encode64(u.Encrypt([]byte(f.Name), UserLog.Key))
		entra = true
	}
	if entra {
		jsonName = aux
	}
	jsonContent := u.Encode64(u.Encrypt(f.Content, UserLog.Key))
	jsonHash := u.Encode64(u.Encrypt(f.HashUser, UserLog.Key))

	data.Set("cmd", cmd)
	data.Set("id", id)
	data.Set("name", jsonName)
	data.Set("content", jsonContent)
	data.Set("hash", jsonHash)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	if resp.Ok {
		fmt.Println(resp.Msg)
	}
	// u.GFicheros = make(map[string]u.Fichero)
	Opciones(resp)
}

///////////////////////////////////////////
///////		DescargarFich	    	///////
///////////////////////////////////////////

func DescargarFich(cmd string, filename string) {
	filepath := "../archivos/"

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("filename", filename)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	var f = u.Fichero{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	content := u.Decrypt(f.Content, UserLog.Key)

	createFile, err := os.Create(filepath + filename)
	u.Chk(err)

	createFile.WriteString(string(content))
	createFile.Close()

	if resp.Ok {
		fmt.Println("Fichero descargado")
	} // copiamos la respuesta a la salida estándar
	Opciones(resp)
}

///////////////////////////////////////////
///////			SubirFich			///////
///////////////////////////////////////////

func SubirFich(cmd string, filename string) {
	filepath := "../archivos/" + filename
	if isExist(filepath) {
		file, err := os.Open(filepath)
		if err != nil {
			panic(err)
		}

		byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero en la variable en bytes
		file.Close()
		f := u.Fichero{}
		json.Marshal(byteValue)
		json.Unmarshal(byteValue, &f.Content)

		f.Name = []byte(filename)
		f.HashUser = UserLog.Key
		f.Content = byteValue

		jsonName := u.Encode64(u.Encrypt([]byte(f.Name), UserLog.Key))
		var aux = string(u.Decode64(jsonName))
		var entra = false
		for ok := true; ok; ok = strings.ContainsAny(aux, "/") {
			aux = u.Encode64(u.Encrypt([]byte(f.Name), UserLog.Key))
			entra = true
		}
		if entra {
			jsonName = aux
		}
		jsonContent := u.Encode64(u.Encrypt(f.Content, UserLog.Key))
		jsonHash := u.Encode64(u.Encrypt(f.HashUser, UserLog.Key))

		data := url.Values{}
		data.Set("cmd", "subirFichero")
		data.Set("id", filename)
		data.Set("name", jsonName)
		data.Set("content", jsonContent)
		data.Set("hash", jsonHash)

		r, err := u.Client.PostForm("https://localhost:10443", data)
		u.Chk(err)
		resp := u.Resp{}
		byteResp, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		json.Unmarshal(byteResp, &resp)
		fmt.Println(resp.Msg)
		if resp.Ok {
			err = os.Remove(filepath)
			u.Chk(err)
		} // copiamos la respuesta a la salida estándar
		Opciones(resp)
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

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("filename", filename)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	var f = u.Fichero{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	hash := u.Decrypt(f.HashUser, UserLog.Key)
	content := u.Decrypt(f.Content, UserLog.Key)

	if resp.Ok {
		if bytes.Equal(hash, UserLog.Key) {
			fmt.Println("Nombre del fichero: " + filename)
			fmt.Println("Fecha de creación: ", f.Fecha)
			fmt.Println("Contenido: " + string(content))
		} else {
			fmt.Println("ERROR: No tiene permisos para leer el fichero")
		}
	}
	Opciones(resp)
}

///////////////////////////////////////////
///////			ListarFICH			///////
///////////////////////////////////////////

func ListarFich(cmd string) {
	fmt.Println("Leyendo fichero ...")
	fmt.Println("------------------------")

	data := url.Values{}
	data.Set("cmd", cmd)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	f := u.FicherosRegistrados{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	ficheros := f.Ficheros

	if resp.Ok {
		fmt.Println("Ficheros de " + UserLog.Name)
		for fich := range ficheros {
			hash := u.Decrypt(ficheros[fich].HashUser, UserLog.Key)
			if bytes.Equal(hash, UserLog.Key) {
				fmt.Println("- " + fich)
			}
		}
	}
	Opciones(resp)
}

///////////////////////////////////////////
///////			EliminarFICH		///////
///////////////////////////////////////////

func EliminarFich(cmd string, filename string) {

	data := url.Values{}
	data.Set("cmd", "leerFichero")
	data.Set("filename", filename)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	var f = u.Fichero{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	hash := u.Decrypt(f.HashUser, UserLog.Key)

	if resp.Ok {
		if bytes.Equal(hash, UserLog.Key) {
			dataRemove := url.Values{}
			dataRemove.Set("cmd", cmd)
			dataRemove.Set("filename", filename)

			r, err := u.Client.PostForm("https://localhost:10443", dataRemove)
			u.Chk(err)

			respRemove := u.Resp{}
			byteValue, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			json.Unmarshal(byteValue, &respRemove)

			fmt.Println(respRemove.Msg)
		} else {
			fmt.Println("ERROR: No tiene permisos para leer el fichero")
		}
	}
	Opciones(resp)
}

func LeerComentarios(cmd string, filename string) {
	data := url.Values{}
	data.Set("cmd", "leerFichero")
	data.Set("filename", filename)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	var f = u.Fichero{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	hash := u.Decrypt(f.HashUser, UserLog.Key)
	if resp.Ok && bytes.Equal(hash, UserLog.Key) {

		dataC := url.Values{}
		dataC.Set("cmd", cmd)
		dataC.Set("filename", filename)
		r, err := u.Client.PostForm("https://localhost:10443", data)
		u.Chk(err)

		respC := u.Resp{}
		byteValueC, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(byteValueC, &respC)
		var f = u.Fichero{}
		json.Unmarshal(u.Decode64(respC.Msg), &f)
		c := f.Comentarios

		fmt.Println("COMENTARIOS: ")
		for k := range c {
			fmt.Println("-------------------------")
			fmt.Println("Comentario")
			fmt.Println("Fecha: ", c[k].Fecha)
			fmt.Println("Mensaje: ", string(u.Decrypt(c[k].Message, UserLog.Key)))
			fmt.Println("-------------------------")
		}
	} else {
		fmt.Println("ERROR: No tiene permisos para leer el fichero")
	}

	Opciones(resp)
}

func AgregarComentarios(cmd string, filename string) {
	data := url.Values{}
	data.Set("cmd", "leerFichero")
	data.Set("filename", filename)

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(byteValue, &resp)
	var f = u.Fichero{}
	json.Unmarshal(u.Decode64(resp.Msg), &f)

	hash := u.Decrypt(f.HashUser, UserLog.Key)
	if resp.Ok && bytes.Equal(hash, UserLog.Key) {
		fmt.Print("Introduzca el comentario:  ")
		comentario := u.LeerTerminal()

		jsonComent := u.Encode64(u.Encrypt([]byte(comentario), UserLog.Key))

		dataC := url.Values{}
		dataC.Set("cmd", cmd)
		dataC.Set("filename", filename)
		dataC.Set("coment", jsonComent)

		r, err := u.Client.PostForm("https://localhost:10443", dataC)
		u.Chk(err)

		respC := u.Resp{}
		byteValue, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(byteValue, &respC)
		if resp.Ok {
			fmt.Println(respC.Msg)
		} // copiamos la respuesta a la salida estándar
	} else {
		fmt.Println("No tiene permisos para comentar este archivo")
	}

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
		fmt.Println("\n---- MENÚ PRINCIPAL ----")
		fmt.Println("1. Crear archivo")
		fmt.Println("2. Listar mis archivos")
		fmt.Println("3. Leer archivo")
		fmt.Println("4. Subir archivo")
		fmt.Println("5. Descargar archivo")
		fmt.Println("6. Eliminar archivo")
		fmt.Println("7. Leer comentarios")
		fmt.Println("8. Agregar comentario")
		fmt.Println("9. Cerrar sesión")
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
			fmt.Println("Se ha seleccionado LISTAR MIS ARCHIVOS")
			fmt.Println("--------------------------------")
			ListarFich("listarFicheros")
		case 3:
			fmt.Println("Se ha seleccionado LEER ARCHIVO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desea ver: ")
			filename := u.LeerTerminal()
			LeerFich("leerFichero", filename)
		case 4:
			fmt.Println("Se ha seleccionado SUBIR ARCHIVO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desea subir: ")
			filename := u.LeerTerminal()

			SubirFich("subirFichero", filename)
		case 5:
			fmt.Println("Se ha seleccionado DESCARGAR ARCHIVO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desea descargar: ")
			filename := u.LeerTerminal()

			DescargarFich("descargarFichero", filename)
		case 6:
			fmt.Println("Se ha seleccionado ELIMINAR ARCHIVO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desea eliminar: ")
			filename := u.LeerTerminal()

			EliminarFich("eliminarFichero", filename)
		case 7:
			fmt.Println("Se ha seleccionado LEER COMENTARIO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desee leer: ")
			filename := u.LeerTerminal()

			LeerComentarios("leerComentarios", filename)
		case 8:
			fmt.Println("Se ha seleccionado AGREGAR COMENTARIO")
			fmt.Println("--------------------------------")
			fmt.Print("Introduzca el nombre del fichero que desee comentar: ")
			filename := u.LeerTerminal()

			AgregarComentarios("agregarComentarios", filename)
		case 9:
			fmt.Println("\n¡Hasta luego!")
			return
		default:
			fmt.Println("No es una opción válida introduzca un número entre 1 y 3:")
			Opciones(resp)
			return
		}

	}

}
