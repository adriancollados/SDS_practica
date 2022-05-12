package main

import (
	"fmt"
	"net/http"
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
	default:
		panic("Opcion no encontrada")
	}
}

func main() {

	u.Gusers = make(map[string]u.User)

	fmt.Println("------------------------------------")
	fmt.Println("Bienvenido al sistema de SDS")
	fmt.Println("------------------------------------")
	fmt.Println("El SERVIDOR se ha iniciado...")
	// fmt.Print("Dime la contraseña del servidor: ")
	// key := leerTerminal()
	// data := sha512.Sum512([]byte(key))
	// codee = data[:32] //El codigo es los primeros 32
	// abrirArchivo()
	http.HandleFunc("/", handler)
	u.Chk(http.ListenAndServeTLS(":10443", "../localhost.crt", "../localhost.key", nil))
}
