# Manual usuario </center>

### Descripcion
Se implemento una aplicacion que se despliega en un cluster de kubernetes, las rutas definidas se desplegaran usando el broker llamado Kafka.
  
### Tecnologías utilizadas:
* Postman

Se utilizo postman para enviar peticiones POST a la siguiente direccion:
```
http://34.134.29.166:2000/StarGame
```
Esta peticion solicita 2 parametros, la cantidad de jugadores y el id de juego (los valores de dichos parametros deberan de enviarse como string usando comillas dobles).

Le agregamos al body de la peticion la siguiente informacion:

<img width="243" alt="image" src="https://user-images.githubusercontent.com/69278553/166118035-a87b4fb0-f418-4039-b321-10309acb71e5.png">

Y listo podemos enviar peticiones con este formato.

El resultado de la peticion es la siguiente:

<img width="811" alt="image" src="https://user-images.githubusercontent.com/69278553/166117906-e468bc4a-d372-428c-aad7-8d02b4bc3a61.png">

  
