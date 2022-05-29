package fich

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	u "sds/util"
	"time"
)

func AgregarComentarios(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		_, ok := u.GFicheros[req.Form.Get("filename")] // ¿existe ya el fichero?
		if !ok {
			u.Response(w, false, "\nERROR: No existe en la base de datos", nil)
			return
		} else {
			if fileSelected, ok := u.GFicheros[req.Form.Get("filename")]; ok {

				var c u.Comentario
				c.Message = u.Decode64(req.Form.Get("coment"))
				c.Fecha = time.Now()

				length := len(fileSelected.Comentarios)
				fileSelected.Comentarios[length+1] = c
				var code []byte = nil
				Fich := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}
				Fich.Key = u.Codee
				Fich.Ficheros = u.GFicheros
				os.Remove("ficheros.json")
				_, err := os.Create("ficheros.json")
				u.Chk(err)
				jsonF, err := json.Marshal(&Fich)
				u.Chk(err)
				var jsonFD = jsonF
				err = ioutil.WriteFile("ficheros.json", jsonFD, 0644)
				u.Chk(err)
				u.Response(w, true, "Comentario añadido correctamente", u.TokenSesion)
			} else {
				u.Response(w, false, "FALSE", u.TokenSesion)
			}
		}
	}
}
