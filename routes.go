package main

import (
	"github.com/gin-gonic/gin"
	"tradeservice/handlers"
)

func routes(r *gin.Engine) {

	root := r.Group("/v1")
	{
		offers := root.Group("/offers")
		{
			offers.POST("", handlers.CreateOffer)
			offers.GET("/:id/execute", handlers.Execute)
		}
		trades := root.Group("/trades")
		{
			trades.GET("", handlers.ListTrades)
		}
	}
}
