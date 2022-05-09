package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	s "sds/cliente/signs"

	"github.com/alexflint/go-arg"
)

//Cliente global
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client = &http.Client{Transport: tr}

func main() {
	var args struct {
		Operation string `arg:"positional, required" help:"(signup|signin)"`
	}
	fmt.Println("**")
	fmt.Println("** Bienvenido al sistema de Foros de la asignatura de SDS en 20/21 ")
	fmt.Println("**")

	parser := arg.MustParse(&args)

	switch args.Operation {
	case "signup":
		s.Signup(client, "signup")
	case "signin":
		s.Signin(client, "signin")
	case "help":
		parser.WriteHelp(os.Stdin)
	default:
		parser.Fail(args.Operation)
	}
}
