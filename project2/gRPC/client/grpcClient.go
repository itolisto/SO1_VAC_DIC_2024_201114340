package main

import (
    "context"
    "net/http"
    "fmt"
    "os"
    "log"

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
    serverUrl = fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
    ctx = context.Background()
)

func allGood(c *gin.Context) {
    c.String(http.StatusOK, "Cours REST API Server Ready")
}

func postCourse(c *gin.Context) {
    var courseRecord courserecord

    if courseDataError := c.BindJSON(&courseRecord); courseDataError != nil {
        fmt.Println(courseDataError)
        c.String(http.StatusBadRequest, courseDataError.Error())
        return
    }

    conn, connErr := grpc.NewClient(serverUrl, opts...)

    if connErr != nil {
        log.Fatalf("fail to dial: %v", connErr)
        c.String(http.StatusBadRequest, connErr.Error())
        return
    }

    defer conn.Close()
    
    client := pb.NewCourseClient(conn)

    log.Println("gRPC client connected to server", serverUrl)

   

    response, responseErr := client.PostCourse(ctx, &pb.CourseRecord{ 
        Curso: courseRecord.Curso,
        Facultad: courseRecord.Facultad,
        Carrera: courseRecord.Carrera,
        Region: courseRecord.Region,
    })

    if responseErr != nil {
        log.Fatalf("client.ListFeatures failed: %v", responseErr)
        c.String(http.StatusBadRequest, responseErr.Error())
        return
    }

    c.String(http.StatusOK, response.Response)

}

func main() {
    router := gin.Default()
    router.GET("/", allGood)
    router.POST("/course", postCourse)

    router.Run("localhost:8000")
}