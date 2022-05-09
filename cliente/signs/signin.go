package signs

import (
	"fmt"
	"net/http"
	"net/url"
	u "sds/util"
)

func Signin(client *http.Client, cmd string) {
	fmt.Println("Loggear un usuario")
	fmt.Println("--------------------")

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", cmd)
	r, err := client.PostForm("https://localhost:10443", data)
	u.Chk(err)
	fmt.Println(r)
}
