package fich

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	u "sds/util"
)

func ListarFich(w http.ResponseWriter, req *http.Request) {
	if req != nil {

		file, err := os.Open("ficheros.json")

		if err != nil {
			panic(err)
		}

		byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero en la variable en bytes

		fich := u.FicherosRegistrados{Key: nil, Ficheros: nil}

		json.Marshal(byteValue)
		json.Unmarshal(byteValue, &fich)
		u.Response(w, true, u.Encode64(byteValue), u.TokenSesion)

	} else {
		u.Response(w, false, "\nERROR: No existen ficheros en la base de datos", nil)
		return
	}

}
