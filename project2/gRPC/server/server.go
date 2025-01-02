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

type courseServer struct {
    pb.UnimplementedCourseServer
}

func (s *courseServer) PostCourse(_ context.Context, grade *pb.CourseRecord) (*pb.CourseResponse, error) {
    return &pb.CourseResponse{ Response: "gRPC server received course"}, nil
}

func main() {
    url := fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
    listener, err := net.Listen("tcp", url)
    if err != nil {
        log.Fatalf("failed to listen: %v", url)
    }

    log.Println("gRPC server started at ", url)

    grpcServer := grpc.NewServer()
    serverImpl := &courseServer{}

    pb.RegisterCourseServer(grpcServer, serverImpl)

    grpcServer.Serve(listener)
}