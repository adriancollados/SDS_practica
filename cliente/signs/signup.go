package signs

import (
	"fmt"
	"net/http"
	"net/url"
)

func Signup(client *http.Client, cmd string) {
	fmt.Println("Registrar un usuario")
	fmt.Println("--------------------")
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", cmd)
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)
	fmt.Println(r)
}
