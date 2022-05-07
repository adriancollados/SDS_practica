package sdspractica

import (
	"fmt"
	"os"
	"sdsGrupal/client"
	"sdsGrupal/server"
)

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "srv":
			fmt.Println("Entrando en modo servidor...")
			server.Run()
		case "cli":
			fmt.Println("Entrando en modo cliente...")
			client.Run()
		default:
			fmt.Println("Parámetro '", os.Args[1], "' desconocido. ")
		}
	} else {
		fmt.Println("Ha ocurrido un error")
	}
	fmt.Println("Practica sds")
	fmt.Println("Servidor escuchando")
	server.Run()
	fmt.Println("Cliente conectado")
	client.Run()
}
