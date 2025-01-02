package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/IBM/sarama"

	pb "usac.sopes1/grpc/ProtoBuffer"
)

var (
	transactionIdGenerator int32 = 1
	topic                        = "assignment"
	producersLock          sync.Mutex
	config                 = sarama.NewConfig()
	kafkaBrokerUrl         = fmt.Sprintf("%s:%s", os.Getenv("KAFKA_SERVER_HOST"), os.Getenv("KAFKA_SERVER_PORT"))
)

type courseServer struct {
	pb.UnimplementedCourseServer
}

func (s *courseServer) PostCourse(_ context.Context, grade *pb.CourseRecord) (*pb.CourseResponse, error) {
	producersLock.Lock()
	defer producersLock.Unlock()

	suffix := transactionIdGenerator
	transactionIdGenerator++
	fmt.Println(suffix)

	config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)

	localProducer, err := sarama.NewAsyncProducer([]string{kafkaBrokerUrl}, config)

	if err != nil {
		msg := fmt.Sprintf("gRPC server kafka new producer error: %v", err)
		log.Println(msg)
		return nil, err

	}

	err = localProducer.BeginTxn()

	if err != nil {
		msg := fmt.Sprintf("gRPC server kafka transaction producer error: %v", err)
		log.Println(msg)
		return nil, err
	}

	time.Sleep(50 * time.Microsecond)
	localProducer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(grade.Curso)}

	localProducer.Close()

	return &pb.CourseResponse{Response: "Assigment pushed to kafka Broker"}, nil
}

func main() {
	version, err := sarama.ParseKafkaVersion(sarama.MaxVersion.String())

	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config.Version = version
	config.Producer.Idempotent = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Transaction.Retry.Backoff = 10
	config.Producer.Transaction.ID = "txn_producer"
	config.Net.MaxOpenRequests = 1
	// config.Producer.Return.Successes = true

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
