package main

import (
	"fmt"
	s "sds/cliente/signs"
	m "sds/util"
)

func main() {
	fmt.Println("---------------------------------------")
	fmt.Println("Bienvenido al sistema de SDS (CLIENTE)")
	fmt.Println("---------------------------------------")

	for {
		fmt.Println("\nIntroduzca una opción:")
		fmt.Println("-----------------------")
		fmt.Println("1. Registrar Usuario")
		fmt.Println("2. Iniciar Sesión")
		fmt.Println("-----------------------")
		fmt.Println("3. Salir del programa")
		fmt.Print("\nOpción: ")
		option := m.LeerTerminal()
		if option == "1" || option == "2" || option == "3" {
			switch option {
			case "1":
				s.Signup(m.Client, "signup")
			case "2":
				s.Signin(m.Client, "signin")
			case "3":
				fmt.Println("\n¡Hasta luego!")
				return
			}
		} else {
			fmt.Println("No es una opción válida introduzca un número entre 1 y 3:")
		}

		for {
			fmt.Print("\n¿Desea realizar otra operación? (s/n): ")
			continuar := m.LeerTerminal()

			if continuar != "s" && continuar != "n" {
				fmt.Println("\nPor favor, introduzca una respuesta válida")
			} else if continuar == "n" {
				fmt.Println("\n¡Hasta luego!")
				return
			} else if continuar == "s" {
				break
			}
		}

	}
}
