# SDS_practica
Práctica de SDS de Alejandro Company y Adrián Collados

Cómo ejecutar la práctica:

Para el cliente:

    + cd/cliente
    + go run client.go

Para el servidor:

    + cd/servidor
    + go run server.go

LA CONTRASEÑA DEL SERVIDOR ES: 123

Instrucciones carpeta /archivos:

    + Los archivos subidos al servidor se encuentran en la carpeta /archivos/subidos. Se encripta su nombre y su contenido.
    + Los archivos dentro de /archivos pueden subirse introduciendo su nombre en el programa.
    + Las descargas de archivos se efectúan en /archivos

Instrucciones creado y subida de archivos:

    + En nuestra simulación sólo hemos contemplado la creación de ficheros de texto (.txt)
    + Si creamos un archivo dentro de la aplicación siempre será txt.
    + Admite la subida de cualquier tipo de fichero (dejamos un pdf de prueba).
        ++ La opción de leer fichero en un archivo que no sea de texto tiene un final trágico XD.
    + Se puede comentar cualquier fichero y queda su fecha registrada.