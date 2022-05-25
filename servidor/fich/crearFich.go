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

		desencryp := u.Decode64(req.Form.Get("fich"))
		user := u.Decrypt([]byte(desencryp), []byte(req.Form.Get("user")))

		jsonString := (string(user))

		var f u.Fichero
		if err := json.Unmarshal([]byte(jsonString), &f); err != nil {
			panic(err)
		}

		u.GFicheros[f.Name] = f

		var code []byte = nil
		Fich := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}
		Fich.Key = u.Codee
		Fich.Ficheros = u.GFicheros
		os.Remove("ficheros.json")
		_, err := os.Create("ficheros.json")
		u.Chk(err)
		jsonF, err := json.Marshal(&Fich)
		u.Chk(err)

		createFile, err := os.Create("../archivos/subidos" + f.Name)
		u.Chk(err)
		defer createFile.Close()

		jsonData, err := json.Marshal(&f.Content)
		u.Chk(err)
		jsonData = []byte(u.Encode64(u.Encrypt(jsonData, Fich.Key)))

		createFile.WriteString(string(jsonData))

		//Encriptamos el json de los ficheros con el codigo de la contraseña del server
		var jsonFD = jsonF

		err = ioutil.WriteFile("ficheros.json", jsonFD, 0644)
		u.Chk(err)
		u.Response(w, true, "Fichero creado correctamente", u.TokenSesion)
	}
}
