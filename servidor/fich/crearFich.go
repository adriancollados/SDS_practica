package fich

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	u "sds/util"
)

func CrearFich(w http.ResponseWriter, req *http.Request) {

	if req != nil {
		_, ok := u.GFicheros[req.Form.Get("id")] // ¿existe ya el usuario?
		if ok {
			u.Response(w, false, "\nERROR: Ya existe en la base de datos", nil)
			return
		} else {
			var f u.Fichero

			f.Name = u.Decode64(req.Form.Get("name"))
			f.HashUser = u.Decode64(req.Form.Get("hash"))
			f.Content = u.Decode64(req.Form.Get("content"))
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
