# Proyecto 1

## Información de estudiante

- Edgar Mauricio Gómez Flores
- 2011-14340

## Descripción

A continuación una descripción para levantar el proyecto.

### Agente

Para construir los procesos de RAM y CPU, se utiliza el código fuente escrito en `C`, después de compilar, obtendrá un archivo `.ko`, estos archivos luego se instalan.

```bash
$ cd ./agent/cpu
$ make all
$ sudo insmod cpu.ko
```

```bash
$ cd ./agent/ram
$ make all
$ sudo insmod cpu.ko
```

Nota: Revise los registros con `dmesg`

```bash
cat /proc/<nombre_modulo>
```

### Base de datos

Se utiliza una imagen de Docker para utilizar la base de datos. La instancia de MySQL necesita variables de entorno para su configuración inicial, las cuales son:

- `MYSQL_ROOT_PASSWORD` = monitor
- `MYSQL_DATABASE` = monitor
- `MYSQL_USER` = monitor
- `MYSQL_PASSWORD` = monitor

Además el script de inicialización se puede encontrar [aquí](./db/init.sql)

### Servidor

La REST API utiliza NodeJS para exponer tres endpoints que ayudan al usuario a realizar un check de salud para el autoscaling y dos para insertar la información que se recibe de los agentes.

- `Ruta`: `/health`
- `Tipo`: `GET`

- `Ruta`: `/ram`
- `Tipo`: `POST`
- `Cuerpo`:
    - total_ram - Número
    - free_ram - Número
    - used_ram - Número
    - percentage_used - Número

- `Ruta`: `/cpu`
- `Tipo`: `POST`
- `Cuerpo`:
    - percentage_used - Número

### Grafana

Se levanta utilizando el script de Docker Compose [aquí](./docker-compose.yml)

### GCP

Crear una Plantilla de instancia para el componente elástico. En `Compute Engine > Instance Templates > CREATE INSTANCE TEMPLATE`:

- Nombre: "monitorvm"
- Ubicación: "us-central1"
- Tipo de máquina: "E2"
- Disco de arranque: "Ubuntu 20.04"
- Permitir HTTP/HTTPS
- Etiquetas de red: "allin", "allout"
- En `Management > Automation > Startup script` pegue el script para descargar y ejecutar los contenedores del agente.

Utilizar el siguiente script:

```bash
sudo apt-get update && sudo apt-get install ca-certificates curl -y
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
git clone https://github.com/itolisto/SO1_VAC_DIC_2024_201114340.git
sudo apt-get install linux-headers-generic -y
sudo apt-get install make
sudo apt-get install gcc -y
sudo apt install --reinstall gcc-12 -y
cd ./SO1_VAC_DIC_2024_201114340/Proyecto1/agent/cpu
make all
sudo insmod cpu_201114340.ko
cd ../ram
make all
sudo insmod ram_201114340.ko
cd ..
sudo docker compose up
```