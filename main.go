package main

import (
	"fmt"
	"github.com/DZDomi/tradeservice/clients"
	"github.com/DZDomi/tradeservice/models"
	"github.com/gin-gonic/gin"
)

func main() {

	models.InitDB("tradeservice:password@/tradeservice?charset=utf8&parseTime=True&loc=Local", false)
	defer models.DB.Close()

	clients.InitRedis("trade")
	clients.InitLock()

	clients.InitKafka()

	gin.SetMode("debug")
	r := gin.Default()

	routes(r)

	_ = r.Run(fmt.Sprintf(":%d", 8080))
}
