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

type graderecord struct {
  Carnet  string  `json:"carnet"`
  Nombre  string `json:"nombre"`
  Curso   string `json:"curso"`
  Nota    int32  `json:"nota"`
  Semestre string `json:"semestre"`
  Año     int32  `json:"año"`
}

var (
    opts = []grpc.DialOption{ grpc.WithTransportCredentials(insecure.NewCredentials()) }
    serverUrl = fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
)

func allGood(c *gin.Context) {
    c.String(http.StatusOK, "Grade REST API Server Ready")
}

func postGrade(c *gin.Context) {
    var gradeRecord graderecord

    if gradeDataError := c.BindJSON(&gradeRecord); gradeDataError != nil {
        fmt.Println(gradeDataError)
        c.String(http.StatusBadRequest, gradeDataError.Error())
        return
    }

    conn, connErr := grpc.NewClient(serverUrl, opts...)

    if connErr != nil {
        log.Fatalf("fail to dial: %v", connErr)
        c.String(http.StatusBadRequest, connErr.Error())
        return
    }

    defer conn.Close()
    
    client := pb.NewGradeClient(conn)

    log.Println("gRPC client connected to server", serverUrl)

   

    response, responseErr := client.PostGrade(context.Background(), &pb.GradeRecord{ 
        Carnet: gradeRecord.Carnet,
        Nombre: gradeRecord.Nombre,
        Curso: gradeRecord.Curso,
        Nota: gradeRecord.Nota,
        Semestre: gradeRecord.Semestre,
        Year: gradeRecord.Año,
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
    router.POST("/grade", postGrade)

    router.Run("localhost:8000")
}