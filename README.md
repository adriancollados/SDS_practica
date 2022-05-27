# SDS_practica
Práctica grupal de SDS de Alejandro Company y Adrián Collados


A tener en cuenta:

- Controlar que no se puedan meter usuarios con contraseña y usuario vacios

- Controlar que no se puedan meter ficheros vacios

- Revisar petada de cierre de sesion despues de leer o crear un archivo

- Listar archivos

- Codificar nombres de ficheros (recorremos el array de ficheros y los que concidan con el HashUser del usuario iniciado se listan)

- Implementar subida, edicion y borrado de ficheros



-----
Cuando creamos el archivo. 
-----
*** Importante ***
No pasamos la clave del UserLog.Key para el conocimiento 0, es uno de los requisitos opcionales que pide el profe 
*** ----- ***

Desde el cliente le pasamos los campos uno a uno encriptados con la contraseña del cliente. 

En el servidor, comprobamos que ese fichero no existe ya en el map de ficheros, si existe manda un error, si no existe continua con la creación. Recuperamos los datos uno a uno y los encriptamos con la contraseña del server y ya lo guardamos. 


-----
Leer ficheros
-----
