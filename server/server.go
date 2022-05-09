package server

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sdsGrupal/server/models"
	"sdsGrupal/util"
	"time"

	"golang.org/x/crypto/argon2"
)

func comprueba(e error) {
	if e != nil {
		panic(e)
	}
}

var usuarios map[string]models.Usuario

func respuesta(w io.Writer, correcto bool, mensaje string, token []byte) {
	r := models.Resp{Correcto: correcto, Mensaje: mensaje, Token: token}
	rJSON, err := json.Marshal(&r)
	comprueba(err)
	w.Write(rJSON)
}

func Run() {
	usuarios = make(map[string]models.Usuario) //inicializamos el mapa de usuarios

	http.HandleFunc("/", handler)

	comprueba(http.ListenAndServeTLS(":9090", "localhost.crt", "localhost.key", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Servidor conectado")
	fmt.Println("------------------------------------------")
	req.ParseForm()
	w.Header().Set("Content-type", "text/plain")

	switch req.Form.Get("cmd") {
	case "registro":
		_, existe := usuarios[req.Form.Get("usuario")] //comprobamos si el usuario ya existe
		if existe {
			respuesta(w, false, "Usuario ya registrado", nil)
			return
		}
		u := models.Usuario{}
		u.Nombre = req.Form.Get("usuario")

		u.Sal = make([]byte, 16)
		rand.Read(u.Sal)

		contrasenya := util.Decode64(req.Form.Get("password"))

		u.Hash = argon2.IDKey([]byte(contrasenya), u.Sal, 16384, 8, 1, 32)

		u.Ultconexion = time.Now()
		u.Token = make([]byte, 16)
		rand.Read(u.Token)

		usuarios[u.Nombre] = u
		respuesta(w, true, "Usuario registrado correctamente", u.Token)

	case "login":
		u, existe := usuarios[req.Form.Get("usuario")]
		if !existe {
			respuesta(w, false, "Usuario inexistente", nil)
			return
		}

		contrasenya := util.Decode64(req.Form.Get("password"))
		hash := argon2.IDKey([]byte(contrasenya), u.Sal, 16384, 8, 1, 32)
		if !bytes.Equal(u.Hash, hash) {
			respuesta(w, false, "Credenciales inválidas", nil)
		} else {
			u.Ultconexion = time.Now()
			u.Token = make([]byte, 16)
			rand.Read(u.Token)
			usuarios[u.Nombre] = u
			respuesta(w, true, "Credenciales válidas", u.Token)
		}

	case "data":
		u, existe := usuarios[req.Form.Get("usuario")]
		if !existe {
			respuesta(w, false, "No autorizado", nil)
			return
		} else if (u.Token == nil) || (time.Since(u.Ultconexion).Minutes() > 60) {
			respuesta(w, false, "No autorizado", nil)
			return
		} else if !bytes.EqualFold(u.Token, util.Decode64(req.Form.Get("token"))) {
			respuesta(w, false, "No autorizado", nil)
			return
		}

	default:
		respuesta(w, false, "No existe el comando", nil)
	}
}
