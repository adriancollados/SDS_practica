package main

import (
	"fmt"
	"os"
	s "sds/cliente/signs"
	m "sds/util"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		Operation string `arg:"positional, required" help:"(signup|signin)"`
	}
	fmt.Println("**")
	fmt.Println("** Bienvenido al sistema de SDS en 21/22 ")
	fmt.Println("**")

	parser := arg.MustParse(&args)

	switch args.Operation {
	case "signup":
		s.Signup(m.Client, "signup")
	case "signin":
		s.Signin(m.Client, "signin")
	case "help":
		parser.WriteHelp(os.Stdin)
	default:
		parser.Fail(args.Operation)
	}
}
