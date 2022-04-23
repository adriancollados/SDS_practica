package sdspractica

import (
	"fmt"
	"sdsGrupal/client"
	"sdsGrupal/server"
)

func main() {
	fmt.Println("Practica sds")
	fmt.Println("Servidor escuchando")
	server.Run()
	fmt.Println("Cliente conectado")
	client.Run()
}
