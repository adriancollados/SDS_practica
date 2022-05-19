package fich

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	u "sds/util"
	"time"
)

func Fichup(filename string, req *http.Request) {
	file, err := os.Create(filename) // crea el fichero de destino (servidor)
	u.Chk(err)
	defer file.Close()      // cierra el fichero al salir de ámbito
	t := time.Now()         // timestamp para medir el tiempo
	io.Copy(file, req.Body) // copia desde el Body del request al fichero con streaming

	m := runtime.MemStats{} // obtiene información acerca del uso de memoria
	runtime.ReadMemStats(&m)
	fmt.Println("SRV::", time.Since(t), ":: Fich escrito")  //imprime tiempo
	fmt.Println("SRV:: memoria ", m.TotalAlloc/1024, " KB") // imprime la memoria total

	fmt.Println("¡Fin!") // devuelve un mensaje al cliente
}
