package main

import (
    "context"
    "os"
    "log"
    "net"
    "fmt"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "usac.sopes1/grpc/ProtoBuffer"
)

type gradesServer struct {
    pb.UnimplementedGradeServer
}

func (s *gradesServer) PostGrade(context.Context, *pb.GradeRecord) (*pb.GradeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostGrade not implemented")
}

func main() {
    url := fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
    listener, err := net.Listen("tcp", url)
    if err != nil {
        log.Fatalf("failed to listen: %v", url)
    }

    log.Println("gRPC server started at ", url)


    grpcServer := grpc.NewServer()
    serverImpl := &gradesServer{}

    pb.RegisterGradeServer(grpcServer, serverImpl)

    grpcServer.Serve(listener)
}