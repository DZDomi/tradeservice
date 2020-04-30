package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tradeservice/clients"
	"tradeservice/models"
	"tradeservice/prices"
)

func main() {

	models.InitDB("tradeservice:password@/tradeservice?charset=utf8&parseTime=True&loc=Local", false)
	defer models.DB.Close()

	clients.InitRedis("trade")
	clients.InitLock()

	go prices.Subscribe()

	gin.SetMode("debug")
	r := gin.Default()

	routes(r)
	//println(models.DB.Error)

	_ = r.Run(fmt.Sprintf(":%d", 8080))
}
