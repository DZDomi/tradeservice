package main

import (
	"github.com/gin-gonic/gin"
	"tradeservice/handlers"
)

func routes(r *gin.Engine) {

	root := r.Group("/v1")
	{
		offer := root.Group("/offer")
		{
			offer.POST("", handlers.CreateOffer)
			offer.GET("/:id/execute", handlers.Execute)
		}
	}
}
