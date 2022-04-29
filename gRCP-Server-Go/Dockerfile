# Obtener la imagen que necesitamos
FROM golang

# Asignar el workdir de la api
WORKDIR /

# Copiar todos los archivos al workdir
COPY . .

# Correr 'npm install' para instalar todas las dependencias que necesitemos
RUN go mod download
# Aca pueden existir muchos mas comandos "run"

# Exponer el puerto 4001 al mundo
EXPOSE 50051

# Correr NodeJS
CMD ["go", "run", "gRPC-Server.go"]