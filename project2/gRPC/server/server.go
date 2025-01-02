package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/IBM/sarama"

	pb "usac.sopes1/grpc/ProtoBuffer"
)

var (
	topic = "assignment"

	config          = sarama.NewConfig()
	producer        sarama.AsyncProducer
	kafkaBrokerrUrl = fmt.Sprintf("%s:%s", os.Getenv("KAFKA_SERVER_HOST"), os.Getenv("KAFKA_SERVER_PORT"))
)

type courseServer struct {
	pb.UnimplementedCourseServer
}

func (s *courseServer) PostCourse(_ context.Context, grade *pb.CourseRecord) (*pb.CourseResponse, error) {

	err := producer.BeginTxn()

	if err != nil {
		msg := fmt.Sprintf("gRPC server kafka transaction producer error: %v", err)
		log.Print(msg)
		return &pb.CourseResponse{Response: msg}, nil
	}

	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder("test")}
	err = producer.CommitTxn()

	if err != nil {
		log.Printf("Producer unable to commit transacction %s", err)

		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Printf("Producer: producer is in a fatal state, need to recreate it")
				return &pb.CourseResponse{Response: "Kafka Producer: producer is in a fatal state, need to recreate it"}, nil
			}

			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Printf("Producer: unable to abort transaction: %+v", err)
					continue
				}

				return &pb.CourseResponse{Response: "Producer, unable to commit transaction"}, nil
			}

			// if not you can retry
			err = producer.CommitTxn()

			if err != nil {
				log.Printf("Producer: unable to commit txn %s", err)
				continue
			}
		}
	}

	return &pb.CourseResponse{Response: "Assigment pushed to kafka Broker"}, nil
}

//     var (
//         wg                                  sync.WaitGroup
//         enqueued, successes, producerErrors int
//     )

//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         for range producer.Successes() {
//                 successes++
//         }
//     }()

//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         for err := range producer.Errors() {
//                 log.Println(err)
//                 producerErrors++
//         }
//     }()

//     ProducerLoop:
//     for {
//         message := &ProducerMessage{Topic: "my_topic", Value: StringEncoder("testing 123")}
//         select {
//         case producer.Input() <- message:
//             enqueued++

//         case <-signals:
//             producer.AsyncClose() // Trigger a shutdown of the producer.
//             break ProducerLoop
//         }
//     }

//     wg.Wait()

//     log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)

//     return &pb.CourseResponse{ Response: "gRPC server received course"}, nil
// }

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
	config.Producer.Return.Successes = true

	localProducer, err := sarama.NewAsyncProducer([]string{kafkaBrokerrUrl}, config)

	if err != nil {
		msg := fmt.Sprintf("gRPC server kafka transaction producer error: %v", err)
		log.Println(msg)
		return

	}

	producer = localProducer
	defer producer.AsyncClose()

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
