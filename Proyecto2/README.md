# Descripción
En este proyecto utilizaremos Locust (biblioteca de Python) para enviar tráfico a un controlador de Ingress de Kubernetes, que es un servicio de Kubernetes con Linkerd instalado. Configuraremos 2 rutas usando Linkerd, cada una recibirá el 50% del tráfico y ambas transferirán datos a una base de datos. 

La primera ruta enviará tráfico a un cliente gRPC escrito en Golang que lo redirigirá a un servidor gRPC, también escrito en Golang. Este servidor tendrá una conexión con una base de datos MongoDB y escribirá los datos allí. La segunda ruta estará escrita en Rust; será un servidor conectado tanto a una base de datos Redis como a MongoDB, escribiendo la información recibida en ambas bases de datos. 

Cada una de estas "rutas" será un objeto de despliegue de Kubernetes (Deployment), con un mínimo de 1 réplica y un máximo de 3, y el uso de CPU no debe superar el 50%. Los datos transmitidos serán notas de estudiantes universitarios, y las visualizaremos conectando un servidor Grafana a las bases de datos.

# Instrucciones

## Crear un JSON de ejemplo de Cursos

1. Primero necesitamos generar datos en el siguiente formato usando JSON. Generaremos un 'pool' de entradas ficticias utilizando Python en un archivo llamado `generator.py`:

```json
{
  "curso": "SO1",
  "facultad": "ingenieria",
  "carrera": "sistemas",
  "region": "NORTE"
}
```

2. Inicia un entorno virtual de Python para encapsular las dependencias del proyecto. Consulta [esta](https://docs.python.org/3/library/venv.html#creating-virtual-environments) documentación para los pasos detallados. Básicamente, si usas Python 3.5 o superior, utiliza la aplicación `venv` para crear entornos virtuales.

3. Crea el entorno virtual ejecutando en la terminal: `python -m venv <path_donde_guardar_el_entorno>`. Por ejemplo, si estás en el directorio deseado, usa `python -m venv venv`. Esto crea el directorio y el entorno virtual.

4. Activa tu entorno virtual según la documentación [aqui](https://docs.python.org/3/library/venv.html#how-venvs-work). En Windows, si usas Git Bash, ejecuta `source venv/scripts/activate`. Consulta la combinación de comandos para bash/zsh y cmd.exe/PowerShell en la guía vinculada.

5. Escribe el archivo generador de cursos usando las bibliotecas "json", "random" e "io". Una vez creado el archivo Python llamado "gradesJsonGenerator.py", ejecútalo con `python gradesJsonGenerator.py`. Obtendrás un JSON con notas de ejemplo.

## Configurar Locust

1. Instala Locust siguiendo la [documentación oficial](https://docs.locust.io/en/stable/installation.html). Básicamente, ejecuta: `pip install locust`.

2. Escribe el archivo de prueba de Locust llamado "locustfile.py". Este archivo realizará, en nuestro caso, solicitudes POST al controlador de Ingress. Consulta la [guía rápida](https://docs.locust.a/en/stable/quickstart.html) y para personalizaciones adicionales revisa [aquí](https://docs.locust.io/en/stable/writing-a-locustfile.html#).

3. Ejecuta las tareas de Locust leyendo el JSON generado con las notas de los estudiantes. Si el archivo se llama "locustfile.py", simplemente ejecuta `locust` en la misma carpeta. Si lo nombraste de otra forma o estás en otro directorio, usa `locust -f <ruta_al_archivo_locust>`.

4. Accede a la interfaz web de Locust para monitorear las solicitudes realizadas en http://localhost:8089.

## Configurar Deployments

Cliente y servidor gRPC en Golang junto al servidor REST en Rust (Primer Deployment)
Tanto el servidor Golang como el de Rust actúan como servidores, pero el primero en la cadena también funciona como cliente. El servidor REST API recibe las solicitudes del controlador Ingress y las reenvía al servidor gRPC que está conectado a MongoDB.

### Crear el servidor REST API y cliente gRPC

Utilizaremos la [documentación oficial de gRPC](https://grpc.io/docs/languages/go/basics/) para crear un servicio en Golang.

1. Crea un módulo en Golang en el directorio del servicio con `go mod init grades/rest-service`.

2. Define un endpoint siguiendo la documentación de Gin. El servidor REST API estará en [grpcClient.go](./gRPC/client/grpcClient.go).

3. Instala las dependencias con `go get .`.

4. Configura las variables de entorno necesarias:

```bash
export GRPC_CLIENT_PORT=8000 \
export GRPC_CLIENT_HOST=localhost
```

5. Ejecuta el servidor con `go run .`.

6. Puedes probarlo manualmente en tu línea de comandos con `curl http://localhost:8000/` o el comando que se muestra a continuación. También puedes probarlo con Locust, pero asegúrate de cambiar la dirección IP, el puerto y el endpoint correctos, y luego simplemente ejecuta el comando `locust` con el entorno virtual activado y desde el directorio donde se encuentra tu archivo de Locust.

```bash
curl http://localhost:8000/course \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"curso": "ANP", "facultad": "Ingenieria", "carrera": "Arte", "region": "METROPOLITANA"}'
```

7. Ahora necesitas seguir la documentación oficial de gRPC en la descripción para crear un cliente gRPC en este mismo servidor. Como se indica allí, vamos a generar el código usando Protocol Buffers. Para hacerlo en Golang, necesitas instalar el compilador de Protocol Buffers y un complemento de Go utilizando [esta guía](https://grpc.io/docs/languages/go/quickstart/#prerequisites). Descarga el archivo adecuado para tu arquitectura desde GitHub, como se indica en las instrucciones. Crea un directorio en cualquier ubicación que prefieras y copia el contenido descargado. Luego, añade la carpeta "bin" a la variable de entorno `PATH` (en macOS y Linux, esto se hace en el archivo .bash o .zsh).

8. Crea una nueva variable de entorno, ya sea como variable de usuario o del sistema, llamada "GOBIN" apuntando al directorio donde deseas que se instalen los complementos de Golang, y luego añádelo a tu variable PATH (en mi caso, estoy usando Windows; si estás en macOS o Linux, añádelo a tu archivo .bash o .zsh).

9. Los Protocol Buffers son una forma de definir un servicio y las estructuras de información que un servicio recibirá y devolverá (si corresponde). Instala los complementos ejecutando los siguientes comandos:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Si quieres saber más sobre cómo se generan los servicios y los tipos de datos (llamados "message", que puedes pensar como el equivalente a la palabra clave "class" en Java), consulta [aquí](https://protobuf.dev/programming-guides/proto3/). Cada "rpc" dentro de un "service" en un archivo .proto es básicamente lo que un endpoint es en REST. Además, consulta [aquí](https://protobuf.dev/reference/go/go-generated/#package), donde se explica que debes definir una opción `go_package` en el archivo .proto.

Una vez que todo esté listo, desde el directorio donde creaste el archivo .proto (en este caso, "gRPC/ProtoBuffer"), simplemente ejecuta:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./courses.proto
```

10. Ahora solo necesitas seguir la documentación de gRPC para crear el cliente. Mi implementación se encuentra en [grpcClient.go](./gRPC/client/grpcClient.go).

11. Más adelante, añadiremos algo de lógica para realizar una solicitud HTTP POST a una API REST en Rust.

### Crear Servidor gRPC

Este nodo es un servidor gRPC que recibe las notas y las reenvía a una cola de Kafka. Primero, simplemente vamos a implementar el servidor gRPC siguiendo la documentación mencionada anteriormente. Creé el archivo de implementación en [server.go](./gRPC/server/server.go).

Si deseas probar estos servidores y clientes, ejecútalos desde una terminal diferente cada uno, usando `go run .` desde el directorio donde se encuentra cada archivo. Además, puedes ejecutar el comando curl mostrado anteriormente. También necesitas definir dos variables de entorno. Si estás usando Bash o Zsh, puedes crear variables de entorno temporales ejecutando los siguientes comandos:

```bash
export GRPC_SERVER_PORT=8010
echo $GRPC_SERVER_PORT
export GRPC_SERVER_HOST=localhost
echo $GRPC_SERVER_HOST
```

Regresaremos más adelante para añadir la lógica necesaria para enviar la información de los cursos a una cola de Kafka.

### Servidor Rust/Cliente Redis

Este nodo recibirá cursos utilizando una API REST creada con Rust.

1. Como estoy usando Windows, consulté [esta documentación](https://learn.microsoft.com/en-us/windows/dev-environment/rust/overview#the-pieces-of-the-rust-development-toolsetecosystem) para familiarizarme con los términos de Rust y [esta otra](https://learn.microsoft.com/en-us/windows/dev-environment/rust/setup) para configurar el entorno de desarrollo para Rust. Básicamente, en Windows necesitas instalar las herramientas de compilación de C++ antes de poder instalar Rust desde su sitio web oficial.

2. Voy a usar el framework "actix web" para crear un servidor web con una API REST siguiendo su [documentación oficial](https://actix.rs/docs/getting-started/), además de usar [JSON](https://actix.rs/docs/extractors#json) para manejar los datos.

3. Crea las variables de entorno utilizando `export`:

```bash
export RUST_SERVER_HOST=localhost
export RUST_SERVER_PORT=8020
```

4. Ejecuta el servidor con `cargo run`.

5. Puedes probarlo usando el siguiente comando:
```bash
curl http://localhost:8020/course \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"curso": "ANP", "facultad": "Ingenieria", "carrera": "Arte", "region": "METROPOLITANA"}'
```

6. Ahora, en el cliente gRPC de la sección anterior, agrega una solicitud HTTP POST hacia este servidor.

### Configurar Cliente Redis en el Servidor de la API REST de Rust

1. Primero, vamos a crear un contenedor Docker para la base de datos Redis utilizando la imagen de Bitnami. Consulta [aquí](https://hub.docker.com/r/bitnami/redis) para obtener más información sobre cómo configurarlo. Básicamente, asumiendo que tienes Docker Desktop instalado y en ejecución, ejecuta el siguiente comando (6379 es el puerto predeterminado). Opcionalmente, puedes usar la versión oficial [redis alpine](https://github.com/docker-library/docs/tree/master/redis) con la imagen `redis:8.0-M02-alpine3.20`. El volumen sería `-v /docker/host/dir:/data`.

```bash
docker run --rm -d -it \
    --name=redis-server \
    -v redis-persistence:/bitnami/redis/data \
    -e REDIS_PASSWORD=course -e REDIS_MASTER_PASSWORD=course \
    -p 6379:6379 \
    bitnami/redis:latest
```

2. Ahora configuraremos el cliente Redis utilizando la biblioteca Rust [actix-extras/actix-session](https://github.com/actix/actix-extras). En el directorio raíz de tu servidor, ejecuta el siguiente comando para agregar la dependencia:
```bash
cargo add actix-session --features=redis-session
```

3. Define las siguientes variables de entorno necesarias para el cliente Redis y otros servicios en tu configuración (por ejemplo, en un archivo `.env` o directamente en tu entorno):

```bash
GRPC_CLIENT_PORT=8000 
GRPC_CLIENT_HOST=<kubernetesObjectTag>

GRPC_SERVER_PORT=8010
GRPC_SERVER_HOST=<kubernetesObjectTag>

RUST_SERVER_PORT=8020
RUST_SERVER_HOST=<kubernetesObjectTag>
RUST_REDIS_PORT=6379
RUST_REDIS_HOST=<kubernetesObjectTag>
```

4. Si deseas configurar las variables de entorno temporalmente en Bash o Zsh, utiliza el siguiente comando:

```bash
export GRPC_CLIENT_PORT=8000 \
export GRPC_CLIENT_HOST=localhost \
export GRPC_SERVER_PORT=8010 \
export GRPC_SERVER_HOST=localhost \
export RUST_SERVER_PORT=8020 \
export RUST_SERVER_HOST=localhost \
export RUST_REDIS_PORT=6379 \
export RUST_REDIS_HOST=localhost
```
