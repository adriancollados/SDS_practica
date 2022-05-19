package signs

import (
	"fmt"
	f "sds/cliente/fich"
	u "sds/util"
)

func Opciones(resp u.Resp) {
	if !resp.Ok {
		fmt.Println("Salir")
		return
	} else {
		for {
			fmt.Println("\n---- MENÚ PRINCIPAL ----")
			fmt.Println("1. Subir archivo")
			fmt.Println("2. Bajar archivo")
			fmt.Println("3. Cerrar sesión")
			fmt.Println("------------------------")
			fmt.Print("¿Qué opción desea realizar? ")
			option := u.StringAInt(u.LeerTerminal())
			fmt.Println("")

			switch option {
			case 1:
				fmt.Println("Se ha seleccionado SUBIR ARCHIVO")
				fmt.Println("--------------------------------")
				fmt.Print("Introduzca el nombre del fichero que desea subir: ")
				filename := u.LeerTerminal()

				f.Fichup(filename)
			case 2:
				fmt.Println("comienza la bajada de archivo")
			case 3:
				fmt.Println("\n¡Hasta luego!")
				return
			default:
				fmt.Println("No es una opción válida introduzca un número entre 1 y 3:")
				Opciones(resp)
				return
			}
		}

	}

}
