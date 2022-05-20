package fich

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"sds/cliente/menu"
	log "sds/cliente/signs"
	u "sds/util"
)

func CrearFich(cmd string) {
	fmt.Println("Creando fichero ...")
	fmt.Println("------------------------")

	if u.GFicheros == nil {
		u.GFicheros = make(map[string]u.Fichero)
	}

	f := u.Fichero{}
	fmt.Println("Nombre del fichero: ")
	f.Name = u.LeerTerminal()
	fmt.Println("Contenido del fichero: ")
	f.Content = []byte(u.LeerTerminal())
	f.HashUser = log.UserLog.Key

	_, ok := u.GFicheros[f.Name]
	if ok {
		fmt.Println("El fichero ya existe")
	} else {
		u.GFicheros[f.Name] = f
	}

	jsonData, err := json.Marshal(&f)
	u.Chk(err)
	jsonData = []byte(u.Encode64(u.Encrypt(jsonData, log.UserLog.Key)))

	data := url.Values{}
	data.Set("cdm", cmd)
	data.Set("fich", string(jsonData))
	data.Set("user", string(log.UserLog.Key))

	r, err := u.Client.PostForm("https://localhost:10443", data)
	u.Chk(err)

	resp := u.Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp)
	u.GFicheros = make(map[string]u.Fichero)
	fmt.Println("Fichero creado")
	fmt.Println("------------------------")
	menu.Opciones(resp)
}
