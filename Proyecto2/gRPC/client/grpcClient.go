package main

import (
    "encoding/json"
    "context"
    "net/http"
    "fmt"
    "bytes"
    "log"
    "os"
    "io"

    "github.com/gin-gonic/gin"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "usac.sopes1/grpc/ProtoBuffer"
)

type courserecord struct {
  Curso    string `json:"curso"`
  Facultad string `json:"facultad"`
  Carrera  string `json:"carrera"`
  Region   string `json:"region"`
}

var (
    opts = []grpc.DialOption{ grpc.WithTransportCredentials(insecure.NewCredentials()) }
    grpcClientUrl = fmt.Sprintf("%s:%s", os.Getenv("GRPC_CLIENT_HOST"), os.Getenv("GRPC_CLIENT_PORT"))
    grpcServerUrl = fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
    ctx = context.Background()
    rustServerUrl = fmt.Sprintf("http://%s:%s", os.Getenv("RUST_SERVER_HOST"), os.Getenv("RUST_SERVER_PORT"))
)

func allGood(c *gin.Context) {
    c.String(http.StatusOK, "Course REST API Server Ready")
}

func postCourse(c *gin.Context) {
    var courseRecord courserecord

    
    if courseDataError := c.BindJSON(&courseRecord); courseDataError != nil {
        fmt.Println(courseDataError)
        c.String(http.StatusBadRequest, courseDataError.Error())
        return
    }

    conn, connErr := grpc.NewClient(grpcServerUrl, opts...)

    if connErr != nil {
        log.Fatalf("fail to dial: %v", connErr)
        c.String(http.StatusBadRequest, connErr.Error())
        return
    }

    defer conn.Close()
    
    client := pb.NewCourseClient(conn)

    log.Println("gRPC client connected to server", grpcServerUrl)

   

    response, responseErr := client.PostCourse(ctx, &pb.CourseRecord{ 
        Curso: courseRecord.Curso,
        Facultad: courseRecord.Facultad,
        Carrera: courseRecord.Carrera,
        Region: courseRecord.Region,
    })

    if responseErr != nil {
        log.Fatalf("gRPC post failed: %v", responseErr)
        c.String(http.StatusBadRequest, responseErr.Error())
        return
    }

    courseJson, courseJsonErr := json.Marshal(courseRecord)

    if courseJsonErr != nil {
        log.Fatalf("parsing back to Json failed: %v", courseJsonErr)
        c.String(http.StatusBadRequest, courseJsonErr.Error())
        return
    }

    rustResponse, rustErr := http.Post(fmt.Sprintf("%s/course", rustServerUrl), "application/json", bytes.NewBuffer(courseJson) )

    if rustErr != nil {
        log.Fatalf("post to Rust server failed: %v", rustErr)
        c.String(http.StatusBadRequest, rustErr.Error())
        return
    }

    defer rustResponse.Body.Close()
    rustBody, rustBodyErr := io.ReadAll(rustResponse.Body)

    if rustBodyErr != nil {
        log.Fatalf("Rust response body failure: %v", rustBodyErr)
        c.String(http.StatusBadRequest, rustBodyErr.Error())
        return
    }

    success := fmt.Sprintf("gRPC server response: %s, Rust REST server response: %s", response.Response, string(rustBody))

    c.String(http.StatusOK, success)

}

func main() {
    router := gin.Default()
    router.GET("/", allGood)
    router.POST("/course", postCourse)

    router.Run(grpcClientUrl)
}