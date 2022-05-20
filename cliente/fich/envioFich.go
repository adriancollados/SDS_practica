package fich

import (
	"fmt"
	"io"
	"net/url"
	"os"
	u "sds/util"
	"time"
)

func Fichup(filename string, cmd string) {
	filepath := "../archivos/" + filename
	if isExist(filepath) {
		fmt.Println("La ruta es correcta")
		// file, err := os.Open(filepath) // abrimos el fichero de origen (cliente)
		// u.Chk(err)
		// defer file.Close() // cerramos al salir de ámbito

		t := time.Now() // timestamp para medir tiempo
		// hacemos un post con formato octet-stream (binario) y ponemos el reader del fichero directamente como Body
		data := url.Values{}
		data.Set("cmd", cmd) // comando (string)
		data.Set("fichero", filename)

		resp, err := u.Client.PostForm("https://localhost:10443", data)
		u.Chk(err)
		fmt.Println("CLIENTE::", time.Since(t), ":: Post realizado") // imprimimos tiempo
		io.Copy(os.Stdout, resp.Body)                                // copiamos la respuesta a la salida estándar
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
