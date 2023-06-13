package main

import (
	"github.com/gin-gonic/gin"
	"id-card-server/service"
	"log"
)

func main() {
	log.Println("Starting")
	r := gin.Default()
	r.POST("/submitData", service.SubmitData)
	r.GET("/createAccount", service.CreateAccount)
	r.POST("/calcData", service.CalcData)
	r.GET("/getAccounts", service.GetAccounts)
	r.GET("/initRedisId/:id", service.InitRedisId)
	r.GET("/test/:id", service.Test)
	r.Run(":7778") // 监听并在 0.0.0.0:8080 上启动服务
}
