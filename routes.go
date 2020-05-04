package main

import (
	"github.com/DZDomi/tradeservice/handlers"
	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {

	root := r.Group("/v1")
	{
		offers := root.Group("/offers")
		{
			offers.POST("", handlers.CreateOffer)
			offers.GET("/:id/accept", handlers.Accept)
		}
		trades := root.Group("/trades")
		{
			trades.GET("", handlers.ListTrades)
		}
	}
}
