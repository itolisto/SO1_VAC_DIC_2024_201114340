# Description
In this project we are going to use Locust(python library) to send traffic to a Kubernetes Ingress controller which
is a Kubernetes service that has Linkerd installed, we will configure 2 routes using Linkerd, each will receive 50% 
of the traffic, both will transfer data to a database. First route will be sending traffic to a gRCP client written 
in Golang that sends it to a gRCP server written in Golang as well, this server has a connection to a Mongo data base
and writes it. Second route is written in Rust, is server that is connected to a redis and the mongo database it writes
the information received to both databases. Each of these 'routes' will be a Kubernetes deployment object, and will have
a minimun of 1 and maximun of 3 replicas and the CPU usage should not be more than 50%. The data that is going to be 
transmitted is collegue students notes, so we will display those notes by connecting a Grafana server to the databases 

# Instructions

## Create Grades Sample Json

1. First we need to generate data in the following format using JSON, we will generate a 'pool' of fake entries using Python in `generator.py`

```json
{
  "carnet": 231565,
  "nombre": "Alumno 1",
  "curso": "SO1",
  "nota": 90,
  "semestre": "2S",
  "year": 2023
}
```

2. We need to start an python virtual environment to encapsulate our projects Python dependencies, so see [this](https://docs.python.org/3/library/venv.html#creating-virtual-environments) documentation to follow the steps. But basically if you're using after Python 3.5
we should use the `venv` application to create these environments.

3. Create the Python environment, run on CLI `python -m venv <path_to_store_environment>` I actually move to the directory I want to store it in and run `python -m venv venv`, this creates the directoy and the venv

4. Then as indicated [here](https://docs.python.org/3/library/venv.html#how-venvs-work) activate your venv, since I'm using git bash in Windows I just run in CLI `source venv/scripts/activate`, if you take a look this is a combination of the bash/zsh and cmd.exe/PowerShell commands

5. We wrote the grades generator file using "json", "random" and "io" libraries. After you create this Python file "gradesJsonGenerator.py" run it using `python gradesJsonGenerator.py`, you will get a json with sample grades

## Set up Locust

1. Install Locust following the [official documentation](https://docs.locust.io/en/stable/installation.html). Basically just do `pip install locust`

2. Write the Locust test in file named "locustfile.py", this code will be basically be doing, in our case, post requests to the Ingress controller. Follow the [official documentation](https://docs.locust.a/en/stable/quickstart.html) and for further customization [check here](https://docs.locust.io/en/stable/writing-a-locustfile.html#). 

3. Run the locust tasks in the locust file, in our case we are reading the generated json that contains the students grades. You have two options if you actually named your file "locustfile.py" just run command `locust` on the cli in the same directory as the file or if you name it differently or you are running in it from a different directoy run `locust -f <path_to_locust_file>`

3. Now you can check access the Locust server to see the requests that are being made using the [web interface](https://docs.locust.io/en/stable/quickstart.html#locust-s-web-interface) running in `http://localhost:8089`

## Set up Routes

### Route 1, Golang gRCP client and server
Basically both are server but the one in the middle is both, client and REST API server. The one that will receive requests from the ingress controller is both. It is an REST API server because it has an 'endpoint' that recieves the grade from Locust that is sending posts request to the it and is a client because it then forwards the information to the following container which is another gRCP server but this one is connected to the Mongo database. We will be using [gRPC's official documentation](https://grpc.io/docs/languages/go/basics/) to create a service in Golang

#### Create REST API-gRPC Client Server
We'll follow [official documentation]{https://go.dev/doc/tutorial/web-service-gin} tutorial to create the REST API using Golang. We assume you've installed Golang

1. Create Golang module in the directory where your service code will live, run `go mod init grades/rest-service`

2. Create the Endpoint following the documentation, our REST API server is inside gRPC/client/grpcClient.go

3. Run `go get .` to add gin module as a dependency to our module.

4. Run the server with `go run .`

5. You can test it manually on your command line with `curl http://localhost:8000/` or the command below. You can also try it with Locust but make sure to change correct IP address, port and endpoint and then just run command `locust` with the venv activated and from the directoy of you locutsfile
```bash
curl http://localhost:8000/grade \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"carnet": "202364507", "nombre": "Mario Boch", "curso": "BD1", "nota": 100, "semestre": "2S", "año": 2024}'
```

6. Now I need to follow the gRPC official documentation in the description to create a gRPC client in this same server. As indicated there, we are going to generate the code using protocol buffers. To do that in Golang we need to install protocol buffers compiler and a Go plugin using [this guide](https://grpc.io/docs/languages/go/quickstart/#prerequisites). Download the proper architecture file from GitHub as indicated in the instructions, Create a directory wherever you want and copy the downloaded content, now add the "bin" folder to the `PATH` variable(in MacOs and Linux that is your .bash or .zsh file). 

7. Create a new environment variable, either a user or system variable, called "GOBIN" pointing to the directory you want the Golang plugins to be installed in and then add it to your "PATH" variable(I'm using windows, if you are on MacOs or Linux add it to your .bash or .zsh file). 

8. Protocol buffers are a way to define a service and the structures of info that a service will receive and return if any. Install plugins, run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` and `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`. If you want to know more about how services and the data types, called "message"(you can think of the "message" keyword as the "class" keyword in Java), are generated look [here](https://protobuf.dev/programming-guides/proto3/). Each "rpc" inside a "sevice" in a ".proto" file is basically what an endpoint is in REST. Also see [here](https://protobuf.dev/reference/go/go-generated/#package) you need to define a `go_package` option in the ".proto" file. After all this is ready, from the directory where you created the ".proto" file, in our case is "gRPC/ProtoBuffer, just run:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./grades.proto
```

9. Now we just follow the gRPC documentation to create the client. My implementation is in "gRPC/client/grpcClient.go"

### 


GRPC_SERVER_PORT = 8010
GRPC_SERVER_HOST = <kubernetesObjectTag>

