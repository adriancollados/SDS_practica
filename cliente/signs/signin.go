package signs

import (
	"fmt"
	"net/http"
	"net/url"
)

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func Signin(client *http.Client, cmd string) {
	fmt.Println("Loggear un usuario")
	fmt.Println("--------------------")

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", cmd)
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)
	fmt.Println(r)
}
