package fich

import (
	"fmt"
	"io"
	"os"
	u "sds/util"
	"time"
)

func Fichup(filename string) {
	file, err := os.Open(filename) // abrimos el fichero de origen (cliente)
	u.Chk(err)
	defer file.Close() // cerramos al salir de ámbito

	t := time.Now() // timestamp para medir tiempo
	// hacemos un post con formato octet-stream (binario) y ponemos el reader del fichero directamente como Body
	resp, err := u.Client.Post("https://localhost:10443", "application/octet-stream", file)
	u.Chk(err)
	fmt.Println("CLIENTE::", time.Since(t), ":: Post realizado") // imprimimos tiempo
	io.Copy(os.Stdout, resp.Body)                                // copiamos la respuesta a la salida estándar
}
