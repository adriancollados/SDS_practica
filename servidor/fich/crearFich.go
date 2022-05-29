package fich

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	u "sds/util"
	"time"
)

func CrearFich(w http.ResponseWriter, req *http.Request) {

	if req != nil {
		_, ok := u.GFicheros[req.Form.Get("id")] // ¿existe ya el usuario?
		if ok {
			u.Response(w, false, "\nERROR: Ya existe en la base de datos", nil)
			return
		} else {
			var f u.Fichero
			var c u.Comentario
			f.Comentarios = make(map[int]u.Comentario)

			f.Name = u.Decode64(req.Form.Get("name"))
			f.HashUser = u.Decode64(req.Form.Get("hash"))
			f.Content = u.Decode64(req.Form.Get("content"))
			c.Message = u.Decode64(req.Form.Get("comments"))
			f.Fecha = time.Now()

			if !bytes.Equal(c.Message, []byte("")) {
				c.Fecha = time.Now()
				f.Comentarios[0] = c
			}
			u.GFicheros[req.Form.Get("id")] = f

			var code []byte = nil
			Fich := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}
			Fich.Key = u.Codee
			Fich.Ficheros = u.GFicheros
			os.Remove("ficheros.json")
			_, err := os.Create("ficheros.json")
			u.Chk(err)
			jsonF, err := json.Marshal(&Fich)
			u.Chk(err)

			createFile, err := os.Create("../archivos/subidos/" + req.Form.Get("name"))
			u.Chk(err)

			createFile.WriteString(req.Form.Get("content"))
			createFile.Close()
			//Encriptamos el json de los ficheros con el codigo de la contraseña del server
			var jsonFD = jsonF

			err = ioutil.WriteFile("ficheros.json", jsonFD, 0644)
			u.Chk(err)
			u.Response(w, true, "Fichero creado correctamente", u.TokenSesion)
		}
	}
}
