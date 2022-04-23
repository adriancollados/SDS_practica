package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sdsGrupal/server"
	"sdsGrupal/util"
)

func comprueba(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cliente := &http.Client{Transport: tr}

	// hash con SHA512 de la contraseña
	keyCliente := sha512.Sum512([]byte("contraseña del cliente"))
	keyLogin := keyCliente[:32]  // una mitad para el login (256 bits)
	keyData := keyCliente[32:64] // la otra para los datos (256 bits)

	pkCliente, err := rsa.GenerateKey(rand.Reader, 1024)
	comprueba(err)
	pkCliente.Precompute() // aceleramos su uso con un precálculo

	pkJSON, err := json.Marshal(&pkCliente) // codificamos con JSON
	comprueba(err)

	keyPub := pkCliente.Public()          // extraemos la clave pública por separado
	pubJSON, err := json.Marshal(&keyPub) // y codificamos con JSON
	comprueba(err)

	var args struct {
		operacion string   `arg:"positional, required" help:"(registro|login|subir|help)`
		tranferir []string `arg:"positional" help:"(list|((upload|download|remove) <files>...)"`
	}

	switch args.operacion {
	case "registro":
		data := url.Values{}                      // estructura para contener los valores
		data.Set("cmd", "register")               // comando (string)
		data.Set("user", "usuario")               // usuario (string)
		data.Set("pass", util.Encode64(keyLogin)) // "contraseña" a base64

		// comprimimos y codificamos la clave pública
		data.Set("pubkey", util.Encode64(util.Compress(pubJSON)))

		// comprimimos, ciframos y codificamos la clave privada
		data.Set("prikey", util.Encode64(util.Encrypt(util.Compress(pkJSON), keyData)))

		r, err := cliente.PostForm("https://localhost:10443", data) // enviamos por POST
		comprueba(err)
		io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
		r.Body.Close()             // hay que cerrar el reader del body
		fmt.Println()

	case "login":
		data = url.Values{}
		data.Set("cmd", "login")                                   // comando (string)
		data.Set("user", "usuario")                                // usuario (string)
		data.Set("pass", util.Encode64(keyLogin))                  // contraseña (a base64 porque es []byte)
		r, err = cliente.PostForm("https://localhost:10443", data) // enviamos por POST
		comprueba(err)
		resp := server.Resp{}
		json.NewDecoder(r.Body).Decode(&resp) // decodificamos la respuesta para utilizar sus campos más adelante
		fmt.Println(resp)                     // imprimimos por pantalla
		r.Body.Close()

	case "subir":
		//codigo para transferir archivos
	case ""
	}

}
