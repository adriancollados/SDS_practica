package fich

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	u "sds/util"
	"time"
)

func Fichup(w http.ResponseWriter, req *http.Request) {
	filename := req.Form.Get("fichero")
	fmt.Println(filename)
	filepath := "../archivos/" + filename
	file, err := os.Open(filepath) // crea el fichero de destino (servidor)

	//aquí deberiamos guardar los datos del fichero en el struct y añadirlo al gFicheros

	u.Chk(err)
	defer file.Close() // cierra el fichero al salir de ámbito
	t := time.Now()    // timestamp para medir el tiempo
	//io.Copy(file, req.Body) // copia desde el Body del request al fichero con streaming

	m := runtime.MemStats{} // obtiene información acerca del uso de memoria
	runtime.ReadMemStats(&m)
	fmt.Println("SRV::", time.Since(t), ":: Fich escrito")  //imprime tiempo
	fmt.Println("SRV:: memoria ", m.TotalAlloc/1024, " KB") // imprime la memoria total

	u.Response(w, true, "Fichero escrito", nil)
}
