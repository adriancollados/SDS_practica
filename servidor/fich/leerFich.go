package fich

import (
	"encoding/json"
	"fmt"
	"net/http"
	u "sds/util"
)

func LeerFich(w http.ResponseWriter, req *http.Request) {

	if req != nil {
		_, ok := u.GFicheros[req.Form.Get("filename")] // ¿existe ya el usuario?
		if !ok {
			u.Response(w, false, "\nERROR: No existe en la base de datos", nil)
			return
		} else {
			if fileSelected, ok := u.GFicheros[req.Form.Get("filename")]; ok {
				fmt.Println("Archivo encontrado, mandando al cliente ...")

				jsonData, err := json.Marshal(fileSelected)
				u.Chk(err)
				jsonDato := u.Encode64(jsonData)
				u.Response(w, true, jsonDato, u.TokenSesion)

			} else {
				u.Response(w, false, "ERROR: No se ha encontrado el archivo", u.TokenSesion)
			}
		}
	}

}
