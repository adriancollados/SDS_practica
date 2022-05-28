package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	f "sds/servidor/fich"
	s "sds/servidor/signs"
	u "sds/util"
)

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              //Es necesario parsear el usuario
	w.Header().Set("Content-Type", "text/plain") //Cabecera estandar

	switch req.Form.Get("cmd") { //Comprobamos el comando desde el cliente
	case "signup":
		fmt.Println("Se ha seleccionado REGISTRO")
		s.Signup(w, req)
	case "signin":
		fmt.Println("Se ha seleccionado LOGIN")
		s.Signin(w, req)
	case "leerFichero":
		fmt.Println("Se ha seleccionado Leer Fichero")
		f.LeerFich(w, req)
	case "listarFicheros":
		fmt.Println("Se ha seleccionado Listar Ficheros")
		f.ListarFich(w, req)
	case "eliminarFichero":
		fmt.Println("Se ha seleccionado Eliminar Fichero")
		f.EliminarFich(w, req)
	case "subirFichero":
		fmt.Println("Se ha seleccionado Subir Fichero")
		f.CrearFich(w, req)
	case "crearFichero":
		fmt.Println("Se ha seleccionado Crear Fichero")
		f.CrearFich(w, req)
	case "descargarFichero":
		fmt.Println("Se ha seleccionado Descargar Fichero")
		f.LeerFich(w, req)
	default:
		panic("Opcion no encontrada")
	}
}

func cargarDatosUsers() {
	file, err := os.Open("users.json") // abrimos el primer fichero (entrada)

	u.Gusers = make(map[string]u.User) //Inicializamos el mapa de los usuarios

	if err != nil {
		file, err = os.Create("users.json") // abrimos el segundo fichero (salida)
		// inicializamos mapa de usuarios
		if err != nil {
			panic(err)
		}
	} else {
		defer file.Close() //Por último cerramos el fichero

		byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero entero en la variable
		var code []byte = nil                //Creamos una variable de bytes
		Regi := u.UsersRegistrados{Key: code, Users: u.Gusers}

		json.Unmarshal(byteValue, &Regi) //Utilizar esta cuando no esté encriptado el fichero

		/*Comprobamos si la contraseña introducida por la persona que inicializa el servidor es la correcta si coincide con la misma que se ha utilizado para
		Codificar los datos del json, si coinciden el codee y la Key del registro, podemos continuar */
		verdad := bytes.Equal(Regi.Key, u.Codee)
		if verdad == false {
			fmt.Println("La contraseña no es la correcta")
			panic(err)
		} else {
			fmt.Println("La contraseña es correcta, puedes continuar")
		}
	}
}

func cargarDatosFicheros() {
	file, err := os.Open("ficheros.json") // abrimos el primer fichero (entrada)

	u.GFicheros = make(map[string]u.Fichero) //Inicializamos el mapa de los ficheros

	if err != nil {
		file, err = os.Create("ficheros.json") // abrimos el segundo fichero (salida)
		// inicializamos mapa de ficheros
		if err != nil {
			panic(err)
		}
	} else {
		defer file.Close() //Por último cerramos el fichero

		byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero entero en la variable
		var code []byte = nil                //Creamos una variable de bytes
		Regi := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}

		json.Unmarshal(byteValue, &Regi) //Utilizar esta cuando no esté encriptado el fichero
	}
}

func main() {
	fmt.Println("------------------------------------")
	fmt.Println("Bienvenido al sistema de SDS")
	fmt.Println("------------------------------------")
	fmt.Println("El SERVIDOR se ha iniciado...")
	fmt.Print("Dime la contraseña del servidor: ")
	key := u.LeerTerminal()
	data := sha512.Sum512([]byte(key))
	u.Codee = data[:32] //El codigo es los primeros 32

	cargarDatosUsers()
	cargarDatosFicheros()
	http.HandleFunc("/", handler)
	u.Chk(http.ListenAndServeTLS(":10443", "../localhost.crt", "../localhost.key", nil))
}
