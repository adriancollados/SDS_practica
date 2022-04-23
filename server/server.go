package server

import (
	"net/http"
)

func comprueba(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {

	http.HandleFunc("/", handler)

	comprueba(http.ListenAndServeTLS(":9090", ""))
}
