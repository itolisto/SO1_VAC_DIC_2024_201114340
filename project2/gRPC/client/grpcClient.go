package main

import (
    "net/http"
    "fmt"

    "github.com/gin-gonic/gin"

    pb "usac.sopes1/grpc/ProtoBuffer"
)

func allGood(c *gin.Context) {
    c.String(http.StatusOK, "Grade REST API Server Ready")
}

func postGrade(c *gin.Context) {
    _, gradeDataError := c.GetRawData()

    if gradeDataError != nil {
        fmt.Println(gradeDataError)
        c.String(http.StatusBadRequest, gradeDataError.Error())
        return
    }

    c.String(http.StatusOK, "Grade received")
}

func main() {
    router := gin.Default()
    router.GET("/", allGood)
    router.POST("/grade", postGrade)

    router.Run("localhost:8000")
}