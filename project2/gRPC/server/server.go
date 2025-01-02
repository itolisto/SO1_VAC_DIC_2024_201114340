package main

import (
    "context"
    "os"
    "log"
    "net"
    "fmt"

    "google.golang.org/grpc"

    pb "usac.sopes1/grpc/ProtoBuffer"
)

type gradesServer struct {
    pb.UnimplementedGradeServer
}

func (s *gradesServer) PostGrade(_ context.Context, grade *pb.GradeRecord) (*pb.GradeResponse, error) {
    return &pb.GradeResponse{ Response: "gRPC server received grade"}, nil
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