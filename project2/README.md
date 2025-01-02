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

## Set up Locust

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

4. Thes as indicated [here](https://docs.python.org/3/library/venv.html#how-venvs-work) activate your venv, since I'm using git bash in Windows I just run in CLI `source venv/scripts/activate`, if you take a look this is a combination of the bash/zsh and cmd.exe/PowerShell commands

5. Install the packages we will need, do `PIP install 