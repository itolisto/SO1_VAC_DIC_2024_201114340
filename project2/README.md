# Description
In this project we are going to use Locust(python library) to send traffic to a Kubernetes Ingress controller which
is a Kubernetes service that has Linkerd installed, we will configure 2 routes using Linkerd, each will receive 50% 
of the traffic, both will transfer data to a database. First route will be sending traffic to a gRCP client written 
in Golang that sends it to a gRCP server written in Golang as well, this server has a connection to a Mongo data base
and writes it. Second route is written in Rust, is server that is connected to a redis and the mongo database it writes
the information received to both databases. Each of these 'routes' will be a Kubernetes deployment object, and will have
a minimun of 1 and maximun of 3 replicas and the CPU usage should not be more than 50%. The data that is going to be 
transmitted is collegue students notes, so we will display those notes by connecting a Grafana server to the databases 

