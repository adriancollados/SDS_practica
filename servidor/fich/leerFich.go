package fich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	u "sds/util"
)

func LeerFich(w http.ResponseWriter, req *http.Request) {

	var fich = u.Fichero{}

	if req != nil {
		filename := req.Form.Get("filename")
		user := u.Decode64(req.Form.Get("user"))

		if fileSelected, ok := u.GFicheros[filename]; ok {
			json.Marshal(fileSelected)
			json.Unmarshal([]byte(fileSelected.HashUser), &fich.HashUser)
			if bytes.Equal(fileSelected.HashUser, user) {
				json.Unmarshal([]byte(fileSelected.Content), &fich.Content)
				fmt.Println("Archivo encontrado, mandando al cliente ...")

				u.Response(w, true, string(fileSelected.Content), u.TokenSesion)
			} else {
				u.Response(w, false, "ERROR: No tiene permisos para leer el archivo", u.TokenSesion)
			}

		} else {
			u.Response(w, false, "ERROR: No se ha encontrado el archivo", u.TokenSesion)
		}
	}

}
