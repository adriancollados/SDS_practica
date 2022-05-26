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

		var name = string(f.Name)
		//Encriptamos los datos con la clave del server
		jsonName, err := json.Marshal(&f.Name)
		u.Chk(err)
		jsonName = []byte(u.Encode64(u.Encrypt(jsonName, u.Codee)))
		jsonHash, err := json.Marshal(&f.Name)
		u.Chk(err)
		jsonHash = []byte(u.Encode64(u.Encrypt(jsonHash, u.Codee)))
		jsonContent, err := json.Marshal(&f.Name)
		u.Chk(err)
		jsonContent = []byte(u.Encode64(u.Encrypt(jsonContent, u.Codee)))
		f.Name = string(jsonName)
		f.HashUser = jsonHash
		f.Content = jsonContent
		u.GFicheros[name] = f

		var code []byte = nil
		Fich := u.FicherosRegistrados{Key: code, Ficheros: u.GFicheros}
		Fich.Key = u.Codee
		Fich.Ficheros = u.GFicheros
		os.Remove("ficheros.json")
		_, err = os.Create("ficheros.json")
		u.Chk(err)
		jsonF, err := json.Marshal(&Fich)
		u.Chk(err)

		createFile, err := os.Create("../archivos/subidos/" + f.Name)
		u.Chk(err)
		defer createFile.Close()

		createFile.WriteString(string(jsonContent))

		//Encriptamos el json de los ficheros con el codigo de la contraseña del server
		var jsonFD = jsonF

		err = ioutil.WriteFile("ficheros.json", jsonFD, 0644)
		u.Chk(err)
		u.Response(w, true, "Fichero creado correctamente", u.TokenSesion)
	}
}
