package signs

import (
	"fmt"
	u "sds/util"
)

func Opciones(resp u.Resp) {
	if !resp.Ok {
		fmt.Println("Salir")
		return
	} else {
		if resp.Msg == "Usuario registrado" {
			fmt.Println("")
			fmt.Println("Inicia sesión")
			Signin(u.Client, "signin")
		} else if resp.Msg == "Credenciales válidas" || resp.Msg == "Añadido a la base de datos" {
			fmt.Println("")
			fmt.Println("---- MENÚ PRINCIPAL ----")
			fmt.Println("1. Subir archivo")
			fmt.Println("2. Bajar archivo")
			fmt.Println("3. Cerrar el programa")
			fmt.Println("------------------------")
			fmt.Print("¿Qué opción quieres realizar? ")
			number := u.StringAInt(u.LeerTerminal())
			fmt.Println("")

			switch number {
			case 1:
				fmt.Println("comienza la subida de archivo")
				//guardar_tema("crear_tema", resp)
				return
			case 2:
				fmt.Println("comienza la bajada de archivo")
				//crear_tema_privado(resp)
			case 3:
				return
			default:
				Opciones(resp)
				return
			}
		}
	}

}
