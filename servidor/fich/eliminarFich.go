package fich

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	u "sds/util"
)

func EliminarFich(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		f, ok := u.GFicheros[req.Form.Get("filename")] // ¿existe ya el fichero?
		if !ok {
			u.Response(w, false, "\nERROR: No existe en la base de datos", nil)
			return
		} else {

			delete(u.GFicheros, req.Form.Get("filename"))
			var code []byte = nil
			Fich := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}
			Fich.Key = u.Codee
			Fich.Ficheros = u.GFicheros
			os.Remove("ficheros.json")
			_, err := os.Create("ficheros.json")
			u.Chk(err)
			jsonF, err := json.Marshal(&Fich)
			u.Chk(err)

			err = os.Remove("../archivos/subidos/" + u.Encode64(f.Name))
			u.Chk(err)

			//Encriptamos el json de los ficheros con el codigo de la contraseña del server
			var jsonFD = jsonF

			err = ioutil.WriteFile("ficheros.json", jsonFD, 0644)
			u.Chk(err)
			u.Response(w, true, "Fichero eliminado correctamente", u.TokenSesion)
		}
	}
}
